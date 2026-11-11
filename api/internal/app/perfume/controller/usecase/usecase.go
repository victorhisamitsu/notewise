// Package usecase orchestrates Notewise's workflows: it validates and
// normalizes input, calls domain services and builds response DTOs.
package usecase

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/victorhisamitsu/notewise/api/internal/app/perfume/controller/dto"
	"github.com/victorhisamitsu/notewise/api/internal/app/perfume/domain/model"
	"github.com/victorhisamitsu/notewise/api/internal/app/perfume/domain/service"
)

// Discovery groups the read-only discovery workflows.
type Discovery interface {
	Search(ctx context.Context, q string) (dto.SearchResponse, error)
	Perfume(ctx context.Context, id string) (dto.PerfumeResponse, error)
	MatchNotes(ctx context.Context, notes []string, tier string) (dto.NoteMatchResponse, error)
	Families(ctx context.Context) (dto.FamiliesResponse, error)
}

// Personal groups the workflows that depend on the user's context.
type Personal interface {
	Today(ctx context.Context, w model.Weather) (dto.TodayResponse, error)
	Profile(ctx context.Context, req dto.ProfileRequest) (dto.ProfileResponse, error)
}

const maxCollection = 30

type discovery struct {
	catalog     service.Catalog
	attribution string
}

// NewDiscovery creates the discovery use cases.
func NewDiscovery(c service.Catalog, attribution string) Discovery {
	return &discovery{catalog: c, attribution: attribution}
}

func invalid(format string, a ...any) error {
	return fmt.Errorf("%w: "+format, append([]any{model.ErrInvalidInput}, a...)...)
}

func (u *discovery) Search(ctx context.Context, q string) (dto.SearchResponse, error) {
	q = strings.TrimSpace(q)
	if l := len([]rune(q)); l < 2 || l > 80 {
		return dto.SearchResponse{}, invalid("query must have 2 to 80 characters")
	}
	res, total, err := u.catalog.Search(ctx, q)
	if err != nil {
		return dto.SearchResponse{}, err
	}
	return dto.SearchResponse{Query: q, Total: total, Results: res, Attribution: u.attribution}, nil
}

func (u *discovery) Perfume(ctx context.Context, id string) (dto.PerfumeResponse, error) {
	if !model.ValidID(id) {
		return dto.PerfumeResponse{}, invalid("invalid perfume id %q", id)
	}
	page, err := u.catalog.PerfumePage(ctx, id)
	if err != nil {
		return dto.PerfumeResponse{}, err
	}
	return dto.PerfumeResponse{Perfume: page.Perfume, Neighbours: page.Neighbours, Attribution: u.attribution}, nil
}

func (u *discovery) MatchNotes(ctx context.Context, notes []string, tier string) (dto.NoteMatchResponse, error) {
	clean := make([]string, 0, len(notes))
	seen := map[string]bool{}
	for _, n := range notes {
		n = strings.ToLower(strings.TrimSpace(n))
		if n == "" || seen[n] {
			continue
		}
		if len(n) > 40 || strings.ContainsAny(n, ",/?#") {
			return dto.NoteMatchResponse{}, invalid("invalid note %q", n)
		}
		seen[n] = true
		clean = append(clean, n)
	}
	if len(clean) < 1 || len(clean) > 4 {
		return dto.NoteMatchResponse{}, invalid("choose between 1 and 4 notes")
	}
	if tier == "" {
		tier = "any"
	}
	switch tier {
	case "any", "top", "heart", "base":
	default:
		return dto.NoteMatchResponse{}, invalid("tier must be any, top, heart or base")
	}
	res, total, err := u.catalog.MatchNotes(ctx, clean, tier)
	if err != nil {
		return dto.NoteMatchResponse{}, err
	}
	return dto.NoteMatchResponse{Notes: clean, Tier: tier, Total: total, Results: res, Attribution: u.attribution}, nil
}

func (u *discovery) Families(ctx context.Context) (dto.FamiliesResponse, error) {
	f, err := u.catalog.Families(ctx)
	if err != nil {
		return dto.FamiliesResponse{}, err
	}
	return dto.FamiliesResponse{Families: f, Attribution: u.attribution}, nil
}

type personal struct {
	recommender service.Recommender
	profiler    service.Profiler
	attribution string
}

// NewPersonal creates the personal use cases.
func NewPersonal(r service.Recommender, p service.Profiler, attribution string) Personal {
	return &personal{recommender: r, profiler: p, attribution: attribution}
}

func (u *personal) Today(ctx context.Context, w model.Weather) (dto.TodayResponse, error) {
	switch {
	case math.IsNaN(w.TemperatureC) || w.TemperatureC < -60 || w.TemperatureC > 60:
		return dto.TodayResponse{}, invalid("temperature out of range")
	case math.IsNaN(w.Humidity) || w.Humidity < 0 || w.Humidity > 100:
		return dto.TodayResponse{}, invalid("humidity must be between 0 and 100")
	case w.Code < 0 || w.Code > 99:
		return dto.TodayResponse{}, invalid("weather code must be between 0 and 99")
	}
	pick, err := u.recommender.Today(ctx, w)
	if err != nil {
		return dto.TodayResponse{}, err
	}
	return dto.TodayResponse{TodayPick: pick, Attribution: u.attribution}, nil
}

func (u *personal) Profile(ctx context.Context, req dto.ProfileRequest) (dto.ProfileResponse, error) {
	ids := make([]string, 0, len(req.PerfumeIDs))
	seen := map[string]bool{}
	for _, id := range req.PerfumeIDs {
		if seen[id] {
			continue
		}
		if !model.ValidID(id) {
			return dto.ProfileResponse{}, invalid("invalid perfume id %q", id)
		}
		seen[id] = true
		ids = append(ids, id)
	}
	if len(ids) > maxCollection {
		return dto.ProfileResponse{}, invalid("a collection can have at most %d perfumes", maxCollection)
	}
	prof, err := u.profiler.Profile(ctx, ids)
	if err != nil {
		return dto.ProfileResponse{}, err
	}
	return dto.ProfileResponse{CollectionProfile: prof, Attribution: u.attribution}, nil
}
