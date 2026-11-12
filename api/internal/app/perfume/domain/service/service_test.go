package service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/victorhisamitsu/notewise/api/internal/app/perfume/domain/model"
)

func TestClassify(t *testing.T) {
	cases := []struct {
		name     string
		w        model.Weather
		occasion string
		mood     string
	}{
		{"hot and humid", model.Weather{TemperatureC: 31, Humidity: 80, IsDay: true}, "summer-days", "hot-humid"},
		{"hot and dry", model.Weather{TemperatureC: 30, Humidity: 30, IsDay: true}, "summer-days", "hot"},
		{"cold", model.Weather{TemperatureC: 9, IsDay: true}, "cold-weather", "cold"},
		{"snow overrides temperature", model.Weather{TemperatureC: 18, Code: 73, IsDay: true}, "cold-weather", "cold"},
		{"rain", model.Weather{TemperatureC: 20, Code: 63, IsDay: true}, "cold-weather", "rainy"},
		{"mild evening", model.Weather{TemperatureC: 21, IsDay: false}, "date-night", "evening"},
		{"mild day", model.Weather{TemperatureC: 22, IsDay: true}, "everyday", "mild"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := Classify(tc.w)
			if c.Occasion != tc.occasion || c.Mood != tc.mood {
				t.Fatalf("got %s/%s, want %s/%s", c.Occasion, c.Mood, tc.occasion, tc.mood)
			}
			if c.Headline == "" || c.Reason == "" {
				t.Fatal("headline and reason must be set")
			}
		})
	}
}

func TestRank_PreferredFirstAvoidedLastStable(t *testing.T) {
	in := []model.Pick{pick("amber1", "amber-resin"), pick("citrus1", "sunlit-citrus"), pick("wood1", "wood-root"), pick("citrus2", "sunlit-citrus")}
	got := Rank(in, []string{"sunlit-citrus"}, []string{"amber-resin"})
	want := []string{"citrus1", "citrus2", "wood1", "amber1"}
	for i := range want {
		if got[i].ID != want[i] {
			t.Fatalf("order = %v, want %v", ids(got), want)
		}
	}
	if in[0].ID != "amber1" {
		t.Fatal("input slice must not be modified")
	}
}

func ids(p []model.Pick) []string {
	out := make([]string, len(p))
	for i := range p {
		out[i] = p[i].ID
	}
	return out
}

func TestToday_UsesOccasionAndCapsPicks(t *testing.T) {
	f := &fakeProvider{}
	for i := 0; i < 10; i++ {
		f.edit = append(f.edit, pick(fmt.Sprint(i), "wood-root"))
	}
	res, err := NewRecommender(f, discard).Today(context.Background(), model.Weather{TemperatureC: 5, IsDay: true})
	if err != nil {
		t.Fatal(err)
	}
	if f.editOccasion != "cold-weather" || res.Occasion != "cold-weather" {
		t.Fatalf("occasion = %q", f.editOccasion)
	}
	if len(res.Picks) != 6 {
		t.Fatalf("picks = %d, want 6", len(res.Picks))
	}
}

func TestPerfumePage_DegradesWhenNeighboursFail(t *testing.T) {
	f := &fakeProvider{
		perfumes:      map[string]model.Perfume{"a": perfume("a", "Wood & Root", nil, nil, nil)},
		neighboursErr: model.ErrUpstream,
	}
	page, err := NewCatalog(f, discard).PerfumePage(context.Background(), "a")
	if err != nil {
		t.Fatalf("page must be served without neighbours, got %v", err)
	}
	if page.Perfume.ID != "a" || page.Neighbours == nil || len(page.Neighbours) != 0 {
		t.Fatalf("unexpected page: %+v", page)
	}
}

func TestPerfumePage_NotFound(t *testing.T) {
	f := &fakeProvider{perfumes: map[string]model.Perfume{}}
	if _, err := NewCatalog(f, discard).PerfumePage(context.Background(), "x"); !errors.Is(err, model.ErrPerfumeNotFound) {
		t.Fatalf("err = %v", err)
	}
}

func TestNotePage_UsesMatchAndSpacesInSlug(t *testing.T) {
	f := &fakeProvider{notes: map[string]model.Note{"pink-pepper": {ID: "pink-pepper", Name: "Pink Pepper"}}}
	page, err := NewCatalog(f, discard).NotePage(context.Background(), "pink-pepper")
	if err != nil {
		t.Fatal(err)
	}
	if page.Total != 42 || page.Perfumes[0].ID != "m-pink pepper" {
		t.Fatalf("page = %+v (the match must receive the note name, not the slug)", page)
	}
}

func TestNotePage_FallsBackToShowcasedWhenMatchFails(t *testing.T) {
	showcased := []model.NoteMatch{{Summary: model.Summary{ID: "shalimar-guerlain"}}}
	f := &fakeProvider{
		notes:    map[string]model.Note{"vanilla": {ID: "vanilla", ShowcasedIn: showcased}},
		matchErr: model.ErrUpstream,
	}
	page, err := NewCatalog(f, discard).NotePage(context.Background(), "vanilla")
	if err != nil {
		t.Fatalf("note page must degrade, got %v", err)
	}
	if page.Total != 1 || page.Perfumes[0].ID != "shalimar-guerlain" {
		t.Fatalf("page = %+v", page)
	}
}

func TestNotePage_NotFound(t *testing.T) {
	f := &fakeProvider{notes: map[string]model.Note{}}
	if _, err := NewCatalog(f, discard).NotePage(context.Background(), "nothing"); !errors.Is(err, model.ErrNoteNotFound) {
		t.Fatalf("err = %v", err)
	}
}
