package usecase

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/victorhisamitsu/notewise/api/internal/app/perfume/controller/dto"
	"github.com/victorhisamitsu/notewise/api/internal/app/perfume/domain/model"
	"github.com/victorhisamitsu/notewise/api/internal/app/perfume/domain/service"
)

type fakeCatalog struct {
	calls     int
	lastNotes []string
}

func (f *fakeCatalog) Search(context.Context, string) ([]model.Summary, int, error) {
	f.calls++
	return []model.Summary{{ID: "a"}}, 1, nil
}
func (f *fakeCatalog) PerfumePage(context.Context, string) (service.PerfumeWithNeighbours, error) {
	f.calls++
	return service.PerfumeWithNeighbours{}, nil
}
func (f *fakeCatalog) MatchNotes(_ context.Context, n []string, _ string) ([]model.NoteMatch, int, error) {
	f.calls++
	f.lastNotes = n
	return nil, 0, nil
}
func (f *fakeCatalog) Families(context.Context) ([]model.Family, error) { return nil, nil }

type fakeProfiler struct{ ids []string }

func (f *fakeProfiler) Profile(_ context.Context, ids []string) (model.CollectionProfile, error) {
	f.ids = ids
	return model.CollectionProfile{}, nil
}

type fakeRecommender struct{ calls int }

func (f *fakeRecommender) Today(context.Context, model.Weather) (model.TodayPick, error) {
	f.calls++
	return model.TodayPick{}, nil
}

func TestSearch_Validation(t *testing.T) {
	c := &fakeCatalog{}
	u := NewDiscovery(c, "attr")
	for _, q := range []string{"", " a ", strings.Repeat("x", 81)} {
		if _, err := u.Search(context.Background(), q); !errors.Is(err, model.ErrInvalidInput) {
			t.Errorf("q=%q: err = %v", q, err)
		}
	}
	if c.calls != 0 {
		t.Fatal("service must not be called for invalid input")
	}
	res, err := u.Search(context.Background(), "  dior ")
	if err != nil || res.Query != "dior" || res.Attribution != "attr" {
		t.Fatalf("res=%+v err=%v", res, err)
	}
}

func TestPerfume_RejectsInvalidID(t *testing.T) {
	c := &fakeCatalog{}
	if _, err := NewDiscovery(c, "").Perfume(context.Background(), "../secret"); !errors.Is(err, model.ErrInvalidInput) || c.calls != 0 {
		t.Fatalf("err=%v calls=%d", err, c.calls)
	}
}

func TestMatchNotes_NormalizesAndValidates(t *testing.T) {
	c := &fakeCatalog{}
	u := NewDiscovery(c, "")
	res, err := u.MatchNotes(context.Background(), []string{" Vanilla", "vanilla", "", "Tonka Bean"}, "")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(c.lastNotes, "|") != "vanilla|tonka bean" || res.Tier != "any" {
		t.Fatalf("notes=%v tier=%q", c.lastNotes, res.Tier)
	}
	bad := [][]string{{}, {"a", "b", "c", "d", "e"}, {"a,b"}}
	for _, n := range bad {
		if _, err := u.MatchNotes(context.Background(), n, "any"); !errors.Is(err, model.ErrInvalidInput) {
			t.Errorf("notes=%v: err = %v", n, err)
		}
	}
	if _, err := u.MatchNotes(context.Background(), []string{"rose"}, "middle"); !errors.Is(err, model.ErrInvalidInput) {
		t.Fatal("invalid tier must be rejected")
	}
}

func TestToday_ValidatesWeather(t *testing.T) {
	r := &fakeRecommender{}
	u := NewPersonal(r, &fakeProfiler{}, "")
	for _, w := range []model.Weather{{TemperatureC: 99}, {Humidity: 120}, {Code: 150}} {
		if _, err := u.Today(context.Background(), w); !errors.Is(err, model.ErrInvalidInput) {
			t.Errorf("w=%+v: err = %v", w, err)
		}
	}
	if r.calls != 0 {
		t.Fatal("recommender must not be called for invalid weather")
	}
}

func TestProfile_DedupesAndLimits(t *testing.T) {
	p := &fakeProfiler{}
	u := NewPersonal(&fakeRecommender{}, p, "")
	if _, err := u.Profile(context.Background(), dto.ProfileRequest{PerfumeIDs: []string{"a-b", "a-b", "c-d"}}); err != nil {
		t.Fatal(err)
	}
	if strings.Join(p.ids, ",") != "a-b,c-d" {
		t.Fatalf("ids = %v", p.ids)
	}
	many := make([]string, 31)
	for i := range many {
		many[i] = "p" + strings.Repeat("x", i+1)
	}
	if _, err := u.Profile(context.Background(), dto.ProfileRequest{PerfumeIDs: many}); !errors.Is(err, model.ErrInvalidInput) {
		t.Fatalf("err = %v", err)
	}
	if _, err := u.Profile(context.Background(), dto.ProfileRequest{PerfumeIDs: []string{"Bad ID"}}); !errors.Is(err, model.ErrInvalidInput) {
		t.Fatalf("err = %v", err)
	}
}
