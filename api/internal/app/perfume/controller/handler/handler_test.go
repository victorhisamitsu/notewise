package handler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/victorhisamitsu/notewise/api/internal/app/perfume/controller/dto"
	"github.com/victorhisamitsu/notewise/api/internal/app/perfume/domain/model"
	"github.com/victorhisamitsu/notewise/api/internal/crosscutting/httpx"
)

type fakeDiscovery struct{ err error }

func (f fakeDiscovery) Search(context.Context, string) (dto.SearchResponse, error) {
	return dto.SearchResponse{}, f.err
}
func (f fakeDiscovery) Perfume(context.Context, string) (dto.PerfumeResponse, error) {
	return dto.PerfumeResponse{}, f.err
}
func (f fakeDiscovery) MatchNotes(context.Context, []string, string) (dto.NoteMatchResponse, error) {
	return dto.NoteMatchResponse{}, f.err
}
func (f fakeDiscovery) Families(context.Context) (dto.FamiliesResponse, error) {
	return dto.FamiliesResponse{}, f.err
}

type fakePersonal struct {
	err     error
	weather model.Weather
}

func (f *fakePersonal) Today(_ context.Context, w model.Weather) (dto.TodayResponse, error) {
	f.weather = w
	return dto.TodayResponse{}, f.err
}
func (f *fakePersonal) Profile(context.Context, dto.ProfileRequest) (dto.ProfileResponse, error) {
	return dto.ProfileResponse{}, f.err
}

var discard = slog.New(slog.NewTextHandler(io.Discard, nil))

func serve(h *Perfume, method, target, body string) *httptest.ResponseRecorder {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/perfumes/{id}", h.Get)
	mux.HandleFunc("GET /api/v1/recommendations/today", h.Today)
	mux.HandleFunc("POST /api/v1/collection/profile", h.Profile)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(method, target, strings.NewReader(body)))
	return rec
}

func TestErrorMapping(t *testing.T) {
	cases := []struct {
		err  error
		want int
	}{
		{nil, 200},
		{model.ErrInvalidInput, 400},
		{model.ErrPerfumeNotFound, 404},
		{fmt.Errorf("x: %w", model.ErrUpstream), 502},
		{context.DeadlineExceeded, 504},
		{errors.New("boom"), 500},
	}
	for _, tc := range cases {
		h := NewPerfume(fakeDiscovery{err: tc.err}, &fakePersonal{}, discard)
		rec := serve(h, http.MethodGet, "/api/v1/perfumes/eau-sauvage-dior", "")
		if rec.Code != tc.want {
			t.Errorf("err=%v: status = %d, want %d", tc.err, rec.Code, tc.want)
		}
		wantCache := httpx.CacheLong
		if tc.err != nil {
			wantCache = httpx.NoStore
		}
		if got := rec.Header().Get("Cache-Control"); got != wantCache {
			t.Errorf("err=%v: cache-control = %q, want %q", tc.err, got, wantCache)
		}
	}
}

func TestServerErrorsDoNotLeakDetails(t *testing.T) {
	h := NewPerfume(fakeDiscovery{err: errors.New("dial tcp 10.0.0.1: secret detail")}, &fakePersonal{}, discard)
	rec := serve(h, http.MethodGet, "/api/v1/perfumes/a-b", "")
	if strings.Contains(rec.Body.String(), "secret detail") {
		t.Fatal("internal error details leaked to the client")
	}
}

func TestToday_ParsesQuery(t *testing.T) {
	p := &fakePersonal{}
	h := NewPerfume(fakeDiscovery{}, p, discard)
	rec := serve(h, http.MethodGet, "/api/v1/recommendations/today?temp=28.5&humidity=75&code=3&isDay=false", "")
	if rec.Code != 200 {
		t.Fatalf("status = %d", rec.Code)
	}
	if p.weather.TemperatureC != 28.5 || p.weather.Humidity != 75 || p.weather.Code != 3 || p.weather.IsDay {
		t.Fatalf("weather = %+v", p.weather)
	}
	if serve(h, http.MethodGet, "/api/v1/recommendations/today?temp=abc&humidity=1&code=1", "").Code != 400 {
		t.Fatal("non-numeric temperature must be rejected")
	}
}

func TestProfile_RejectsMalformedJSON(t *testing.T) {
	h := NewPerfume(fakeDiscovery{}, &fakePersonal{}, discard)
	if rec := serve(h, http.MethodPost, "/api/v1/collection/profile", "{not json"); rec.Code != 400 {
		t.Fatalf("status = %d", rec.Code)
	}
	if rec := serve(h, http.MethodPost, "/api/v1/collection/profile", `{"perfumeIds":["a-b"]}`); rec.Code != 200 || rec.Header().Get("Cache-Control") != httpx.NoStore {
		t.Fatalf("status = %d cache=%q", rec.Code, rec.Header().Get("Cache-Control"))
	}
}
