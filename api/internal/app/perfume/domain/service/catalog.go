package service

import (
	"context"
	"log/slog"

	"golang.org/x/sync/errgroup"

	"github.com/victorhisamitsu/notewise/api/internal/app/perfume/domain/model"
)

// PerfumeWithNeighbours is a perfume page: profile plus similar perfumes.
type PerfumeWithNeighbours struct {
	Perfume    model.Perfume
	Neighbours []model.Neighbour
}

// Catalog serves search, perfume pages, note matching and families.
type Catalog interface {
	Search(ctx context.Context, q string) ([]model.Summary, int, error)
	PerfumePage(ctx context.Context, id string) (PerfumeWithNeighbours, error)
	MatchNotes(ctx context.Context, notes []string, tier string) ([]model.NoteMatch, int, error)
	Families(ctx context.Context) ([]model.Family, error)
}

type catalog struct {
	provider      PerfumeProvider
	maxNeighbours int
	log           *slog.Logger
}

// NewCatalog creates the catalog service.
func NewCatalog(p PerfumeProvider, log *slog.Logger) Catalog {
	return &catalog{provider: p, maxNeighbours: 8, log: log.With("component", "catalog_service")}
}

func (s *catalog) Search(ctx context.Context, q string) ([]model.Summary, int, error) {
	res, total, err := s.provider.Search(ctx, q)
	if err != nil {
		s.log.ErrorContext(ctx, "search failed", "query", q, "err", err)
		return nil, 0, err
	}
	s.log.InfoContext(ctx, "search", "query", q, "results", len(res), "total", total)
	return res, total, nil
}

// PerfumePage loads the profile and the neighbours concurrently. The
// profile is required; neighbours are a nice-to-have, so a failure there
// is logged and the page is still served (graceful degradation).
func (s *catalog) PerfumePage(ctx context.Context, id string) (PerfumeWithNeighbours, error) {
	var (
		page          PerfumeWithNeighbours
		neighboursErr error
	)
	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		p, err := s.provider.Perfume(gctx, id)
		page.Perfume = p
		return err
	})
	g.Go(func() error {
		n, err := s.provider.Neighbours(gctx, id)
		if err != nil {
			neighboursErr = err
			return nil
		}
		if len(n) > s.maxNeighbours {
			n = n[:s.maxNeighbours]
		}
		page.Neighbours = n
		return nil
	})
	if err := g.Wait(); err != nil {
		s.log.WarnContext(ctx, "perfume page failed", "perfume_id", id, "err", err)
		return PerfumeWithNeighbours{}, err
	}
	if neighboursErr != nil {
		s.log.WarnContext(ctx, "neighbours unavailable, serving page without them", "perfume_id", id, "err", neighboursErr)
	}
	if page.Neighbours == nil {
		page.Neighbours = []model.Neighbour{}
	}
	return page, nil
}

func (s *catalog) MatchNotes(ctx context.Context, notes []string, tier string) ([]model.NoteMatch, int, error) {
	res, total, err := s.provider.MatchNotes(ctx, notes, tier)
	if err != nil {
		s.log.WarnContext(ctx, "note match failed", "notes", notes, "tier", tier, "err", err)
		return nil, 0, err
	}
	s.log.InfoContext(ctx, "note match", "notes", notes, "tier", tier, "total", total)
	return res, total, nil
}

func (s *catalog) Families(ctx context.Context) ([]model.Family, error) {
	f, err := s.provider.Families(ctx)
	if err != nil {
		s.log.ErrorContext(ctx, "families failed", "err", err)
	}
	return f, err
}
