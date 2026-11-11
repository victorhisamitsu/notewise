package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/victorhisamitsu/notewise/api/internal/app/perfume/domain/model"
)

// Profiler analyses a collection of perfumes.
type Profiler interface {
	Profile(ctx context.Context, ids []string) (model.CollectionProfile, error)
}

type profiler struct {
	provider    PerfumeProvider
	concurrency int
	log         *slog.Logger
}

// NewProfiler creates the collection profiler.
func NewProfiler(p PerfumeProvider, log *slog.Logger) Profiler {
	return &profiler{provider: p, concurrency: 5, log: log.With("component", "profile_service")}
}

var (
	allSeasons = []string{"spring", "summer", "autumn", "winter"}
	allHours   = []string{"day", "evening", "night"}
)

// Profile loads every perfume (at most "concurrency" requests in flight,
// to stay polite with the provider), then aggregates families, notes,
// seasons and moments, finds gaps and asks for complementary perfumes.
func (s *profiler) Profile(ctx context.Context, ids []string) (model.CollectionProfile, error) {
	var (
		mu       sync.Mutex
		perfumes = make([]model.Perfume, 0, len(ids))
		missing  []string
	)
	g, gctx := errgroup.WithContext(ctx)
	g.SetLimit(s.concurrency)
	for _, id := range ids {
		g.Go(func() error {
			p, err := s.provider.Perfume(gctx, id)
			mu.Lock()
			defer mu.Unlock()
			switch {
			case errors.Is(err, model.ErrPerfumeNotFound):
				missing = append(missing, id)
				return nil
			case err != nil:
				return err
			}
			perfumes = append(perfumes, p)
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		s.log.ErrorContext(ctx, "profile failed", "size", len(ids), "err", err)
		return model.CollectionProfile{}, err
	}
	// Keep the user's order for deterministic output.
	order := make(map[string]int, len(ids))
	for i, id := range ids {
		order[id] = i
	}
	sort.Slice(perfumes, func(i, j int) bool { return order[perfumes[i].ID] < order[perfumes[j].ID] })

	profile := Aggregate(perfumes)
	sort.Strings(missing)
	profile.Missing = append([]string{}, missing...)

	if len(perfumes) > 0 {
		profile.Suggestions = s.suggestions(ctx, perfumes)
	}
	s.log.InfoContext(ctx, "collection profiled", "size", len(ids), "found", len(perfumes), "missing", len(missing))
	return profile, nil
}

// suggestions asks the provider for a complementary wardrobe and drops
// perfumes already in the collection. Failures only cost the suggestions.
func (s *profiler) suggestions(ctx context.Context, perfumes []model.Perfume) []model.Pick {
	loves := make([]string, 0, 5)
	owned := make(map[string]bool, len(perfumes))
	for _, p := range perfumes {
		owned[p.ID] = true
		if len(loves) < 5 {
			loves = append(loves, p.ID)
		}
	}
	picks, err := s.provider.Wardrobe(ctx, loves, 5)
	if err != nil {
		s.log.WarnContext(ctx, "wardrobe suggestions unavailable", "err", err)
		return []model.Pick{}
	}
	out := make([]model.Pick, 0, len(picks))
	for _, p := range picks {
		if !owned[p.ID] {
			out = append(out, p)
		}
	}
	return out
}

// Aggregate computes the statistics of a collection. It is a pure function,
// which keeps the core logic easy to test.
func Aggregate(perfumes []model.Perfume) model.CollectionProfile {
	families, notes, seasons, hours, words := counter{}, counter{}, counter{}, counter{}, counter{}
	for _, p := range perfumes {
		families.add(p.Family.Name)
		for _, n := range p.AllNotes() {
			notes.add(n)
		}
		for _, v := range p.Seasons {
			seasons.add(v)
		}
		for _, v := range p.Hours {
			hours.add(v)
		}
		for _, v := range p.CharacterWords {
			words.add(v)
		}
	}

	prof := model.CollectionProfile{
		Size:           len(perfumes),
		Families:       families.sorted(0),
		TopNotes:       notes.sorted(10),
		Seasons:        ordered(seasons, allSeasons),
		Hours:          ordered(hours, allHours),
		CharacterWords: words.sorted(8),
		Gaps:           []string{},
		Suggestions:    []model.Pick{},
		Missing:        []string{},
	}
	if len(perfumes) == 0 {
		prof.Headline = "Add perfumes to your shelf to discover your scent profile."
		return prof
	}
	for _, s := range allSeasons {
		if seasons[s] == 0 {
			prof.Gaps = append(prof.Gaps, s)
		}
	}
	for _, h := range allHours {
		if hours[h] == 0 {
			prof.Gaps = append(prof.Gaps, h)
		}
	}
	prof.Headline = headline(prof)
	return prof
}

func headline(p model.CollectionProfile) string {
	if len(p.Families) > 1 && p.Families[0].Count == p.Families[1].Count {
		if len(p.TopNotes) > 0 && p.TopNotes[0].Count > 1 {
			return fmt.Sprintf("An eclectic shelf across %d families, tied together by %s.", len(p.Families), strings.ToLower(p.TopNotes[0].Label))
		}
		return fmt.Sprintf("An eclectic shelf across %d families.", len(p.Families))
	}
	fam := p.Families[0].Label
	if len(p.TopNotes) > 0 && p.TopNotes[0].Count > 1 {
		return fmt.Sprintf("Your shelf leans %s, with a soft spot for %s.", fam, strings.ToLower(p.TopNotes[0].Label))
	}
	return fmt.Sprintf("Your shelf leans %s.", fam)
}

type counter map[string]int

func (c counter) add(k string) {
	if k = strings.TrimSpace(k); k != "" {
		c[k]++
	}
}

// sorted returns counts by frequency (ties alphabetically), capped at limit
// (0 = no cap).
func (c counter) sorted(limit int) []model.Count {
	out := make([]model.Count, 0, len(c))
	for k, v := range c {
		out = append(out, model.Count{Label: k, Count: v})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Count != out[j].Count {
			return out[i].Count > out[j].Count
		}
		return out[i].Label < out[j].Label
	})
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out
}

// ordered returns counts for a fixed label order, including zeros.
func ordered(c counter, labels []string) []model.Count {
	out := make([]model.Count, 0, len(labels))
	for _, l := range labels {
		out = append(out, model.Count{Label: l, Count: c[l]})
	}
	return out
}
