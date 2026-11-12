package scento

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/victorhisamitsu/notewise/api/internal/app/perfume/domain/model"
)

// fixtures were captured from the real API and live in testdata/.
var routes = map[string]string{
	"/perfumes/eau-sauvage-dior":            "detail_pyramid.json",
	"/perfumes/hypnotic-poison-dior":        "detail_flat.json",
	"/perfumes/doesnotexist-xyz":            "detail_notfound.json",
	"/perfumes/eau-sauvage-dior/neighbours": "neighbours.json",
	"/search":                               "search.json",
	"/materials/match":                      "match.json",
	"/families":                             "families.json",
	"/edits/cold-weather":                   "edit.json",
	"/materials/vanilla":                    "material.json",
	"/materials/nothingness":                "material_notfound.json",
}

func newTestClient(t *testing.T, h http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	c := NewClient(srv.URL, 2*time.Second, time.Hour, "notewise-test", slog.New(slog.NewTextHandler(io.Discard, nil)))
	c.retryWait = time.Millisecond
	return c
}

func fixtureServer(t *testing.T, calls *atomic.Int64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if r.URL.Path == "/edits/bad" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":{"code":"invalid_input","message":"occasion must be one of: office"}}`))
			return
		}
		f, ok := routes[r.URL.Path]
		if !ok {
			http.NotFound(w, r)
			return
		}
		b, err := os.ReadFile("testdata/" + f)
		if err != nil {
			t.Fatal(err)
		}
		_, _ = w.Write(b)
	}
}

func TestPerfume_WithPyramid(t *testing.T) {
	var calls atomic.Int64
	c := newTestClient(t, fixtureServer(t, &calls))

	p, err := c.Perfume(context.Background(), "eau-sauvage-dior")
	if err != nil {
		t.Fatal(err)
	}
	if p.ID != "eau-sauvage-dior" || p.House != "Dior" || p.Year != 1966 || p.Perfumer != "Edmond Roudnitska" {
		t.Fatalf("unexpected perfume: %+v", p.Summary)
	}
	if p.Family.ID != "sunlit-citrus" {
		t.Fatalf("family id = %q", p.Family.ID)
	}
	if p.Pyramid == nil || p.Pyramid.Top[0] != "Lemon" || p.Pyramid.Base[0] != "Vetiver" {
		t.Fatalf("pyramid = %+v", p.Pyramid)
	}
	if len(p.Notes) != 5 {
		t.Fatalf("flattened notes = %v", p.Notes)
	}
}

func TestPerfume_FlatNotes(t *testing.T) {
	var calls atomic.Int64
	c := newTestClient(t, fixtureServer(t, &calls))

	p, err := c.Perfume(context.Background(), "hypnotic-poison-dior")
	if err != nil {
		t.Fatal(err)
	}
	if p.Pyramid != nil {
		t.Fatal("expected no pyramid")
	}
	if strings.Join(p.Notes, ",") != "Almond,Jasmine Sambac,Vanilla" {
		t.Fatalf("notes = %v", p.Notes)
	}
}

func TestPerfume_NotFoundIsADomainError(t *testing.T) {
	var calls atomic.Int64
	c := newTestClient(t, fixtureServer(t, &calls))
	if _, err := c.Perfume(context.Background(), "doesnotexist-xyz"); !errors.Is(err, model.ErrPerfumeNotFound) {
		t.Fatalf("err = %v, want ErrPerfumeNotFound", err)
	}
}

func TestSearch_MapsIDsAndCleansNames(t *testing.T) {
	var calls atomic.Int64
	c := newTestClient(t, fixtureServer(t, &calls))

	res, total, err := c.Search(context.Background(), "vanilla")
	if err != nil {
		t.Fatal(err)
	}
	if total == 0 || len(res) == 0 {
		t.Fatal("expected results")
	}
	if res[0].ID != "vanilla-28-kayali" {
		t.Fatalf("id = %q", res[0].ID)
	}
	for _, r := range res {
		if strings.HasSuffix(r.Name, ":") {
			t.Fatalf("name not cleaned: %q", r.Name)
		}
	}
}

func TestNeighboursMatchFamiliesEdit(t *testing.T) {
	var calls atomic.Int64
	c := newTestClient(t, fixtureServer(t, &calls))
	ctx := context.Background()

	if n, err := c.Neighbours(ctx, "eau-sauvage-dior"); err != nil || len(n) == 0 || len(n[0].SharedNotes) == 0 {
		t.Fatalf("neighbours: %v %+v", err, n)
	}
	if m, _, err := c.MatchNotes(ctx, []string{"vanilla", "tonka"}, "any"); err != nil || m[0].Tiers["Vanilla"] != "base" {
		t.Fatalf("match: %v %+v", err, m)
	}
	if f, err := c.Families(ctx); err != nil || len(f) != 14 || f[0].ID != "sunlit-citrus" {
		t.Fatalf("families: %v %d", err, len(f))
	}
	if e, err := c.Edit(ctx, "cold-weather"); err != nil || len(e) == 0 || e[0].Why == "" {
		t.Fatalf("edit: %v %+v", err, e)
	}
}

func TestNote(t *testing.T) {
	var calls atomic.Int64
	c := newTestClient(t, fixtureServer(t, &calls))
	n, err := c.Note(context.Background(), "vanilla")
	if err != nil {
		t.Fatal(err)
	}
	if n.ID != "vanilla" || n.SmellsLike == "" || n.Family.ID != "velvet-sweet" || len(n.ShowcasedIn) == 0 {
		t.Fatalf("note = %+v", n)
	}
	if len(n.Aliases) > 8 {
		t.Fatalf("aliases must be capped, got %d", len(n.Aliases))
	}
	if n.ShowcasedIn[0].Tiers["Vanilla"] == "" {
		t.Fatalf("showcased tier missing: %+v", n.ShowcasedIn[0])
	}
	if _, err := c.Note(context.Background(), "nothingness"); !errors.Is(err, model.ErrNoteNotFound) {
		t.Fatalf("err = %v, want ErrNoteNotFound", err)
	}
}

func TestInvalidInputIsNotRetriedNorCached(t *testing.T) {
	var calls atomic.Int64
	c := newTestClient(t, fixtureServer(t, &calls))
	for i := 0; i < 2; i++ {
		if _, err := c.Edit(context.Background(), "bad"); !errors.Is(err, model.ErrInvalidInput) {
			t.Fatalf("err = %v, want ErrInvalidInput", err)
		}
	}
	if calls.Load() != 2 {
		t.Fatalf("calls = %d, want 2 (no retry, no caching of errors)", calls.Load())
	}
}

func TestResponsesAreCached(t *testing.T) {
	var calls atomic.Int64
	c := newTestClient(t, fixtureServer(t, &calls))
	for i := 0; i < 5; i++ {
		if _, err := c.Families(context.Background()); err != nil {
			t.Fatal(err)
		}
	}
	if calls.Load() != 1 {
		t.Fatalf("upstream calls = %d, want 1", calls.Load())
	}
}

func TestRetriesOnceOnServerError(t *testing.T) {
	var calls atomic.Int64
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if calls.Add(1) == 1 {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		b, _ := os.ReadFile("testdata/families.json")
		_, _ = w.Write(b)
	})
	if _, err := c.Families(context.Background()); err != nil {
		t.Fatalf("expected success after retry, got %v", err)
	}
	if calls.Load() != 2 {
		t.Fatalf("calls = %d, want 2", calls.Load())
	}
}

func TestPersistentFailureIsUpstreamError(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusServiceUnavailable) })
	if _, err := c.Families(context.Background()); !errors.Is(err, model.ErrUpstream) {
		t.Fatalf("err = %v, want ErrUpstream", err)
	}
}
