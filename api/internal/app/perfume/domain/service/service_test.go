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

func TestAggregate(t *testing.T) {
	p := Aggregate([]model.Perfume{
		perfume("a", "Wood & Root", []string{"Vanilla", "Cedar"}, []string{"autumn", "winter"}, []string{"evening"}),
		perfume("b", "Wood & Root", []string{"Vanilla", "Bergamot"}, []string{"winter"}, []string{"day"}),
		perfume("c", "Sunlit Citrus", []string{"Lemon"}, []string{"summer"}, []string{"day"}),
	})
	if p.Size != 3 || p.Families[0].Label != "Wood & Root" || p.Families[0].Count != 2 {
		t.Fatalf("families = %+v", p.Families)
	}
	if p.TopNotes[0].Label != "Vanilla" || p.TopNotes[0].Count != 2 {
		t.Fatalf("top notes = %+v", p.TopNotes)
	}
	if len(p.Seasons) != 4 || p.Seasons[3].Label != "winter" || p.Seasons[3].Count != 2 {
		t.Fatalf("seasons = %+v", p.Seasons)
	}
	if fmt.Sprint(p.Gaps) != "[spring night]" {
		t.Fatalf("gaps = %v", p.Gaps)
	}
	if p.Headline != "Your shelf leans Wood & Root, with a soft spot for vanilla." {
		t.Fatalf("headline = %q", p.Headline)
	}
}

func TestAggregate_EclecticHeadline(t *testing.T) {
	p := Aggregate([]model.Perfume{
		perfume("a", "Wood & Root", []string{"Vanilla"}, nil, nil),
		perfume("b", "Sunlit Citrus", []string{"Vanilla"}, nil, nil),
	})
	if p.Headline != "An eclectic shelf across 2 families, tied together by vanilla." {
		t.Fatalf("headline = %q", p.Headline)
	}
}

func TestAggregate_Empty(t *testing.T) {
	p := Aggregate(nil)
	if p.Size != 0 || p.Headline == "" || p.Gaps == nil || p.Families == nil {
		t.Fatalf("empty profile must be fully initialised: %+v", p)
	}
}

func TestProfile_ReportsMissingBoundsConcurrencyAndFiltersOwned(t *testing.T) {
	f := &fakeProvider{perfumes: map[string]model.Perfume{}}
	var ids []string
	for i := 0; i < 12; i++ {
		id := fmt.Sprintf("p%d", i)
		ids = append(ids, id)
		f.perfumes[id] = perfume(id, "Green Leaf", []string{"Mint"}, []string{"spring"}, []string{"day"})
	}
	ids = append(ids, "ghost")
	f.wardrobe = []model.Pick{pick("p0", "green-leaf"), pick("new", "wood-root")}

	prof, err := NewProfiler(f, discard).Profile(context.Background(), ids)
	if err != nil {
		t.Fatal(err)
	}
	if prof.Size != 12 || fmt.Sprint(prof.Missing) != "[ghost]" {
		t.Fatalf("size=%d missing=%v", prof.Size, prof.Missing)
	}
	if f.maxInFlight > 5 {
		t.Fatalf("max in-flight = %d, want <= 5", f.maxInFlight)
	}
	if len(prof.Suggestions) != 1 || prof.Suggestions[0].ID != "new" {
		t.Fatalf("suggestions = %+v (owned perfumes must be filtered)", prof.Suggestions)
	}
}

func TestProfile_UpstreamErrorFails(t *testing.T) {
	f := &fakeProvider{perfumeErr: model.ErrUpstream}
	if _, err := NewProfiler(f, discard).Profile(context.Background(), []string{"a"}); !errors.Is(err, model.ErrUpstream) {
		t.Fatalf("err = %v", err)
	}
}
