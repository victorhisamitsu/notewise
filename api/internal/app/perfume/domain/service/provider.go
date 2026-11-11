// Package service holds Notewise's business rules.
package service

import (
	"context"

	"github.com/victorhisamitsu/notewise/api/internal/app/perfume/domain/model"
)

// PerfumeProvider is the perfume data source (implemented by the Scento
// repository). It is declared here, where it is consumed, so the domain
// does not depend on any concrete API.
type PerfumeProvider interface {
	Search(ctx context.Context, q string) ([]model.Summary, int, error)
	Perfume(ctx context.Context, id string) (model.Perfume, error)
	Neighbours(ctx context.Context, id string) ([]model.Neighbour, error)
	MatchNotes(ctx context.Context, notes []string, tier string) ([]model.NoteMatch, int, error)
	Families(ctx context.Context) ([]model.Family, error)
	Edit(ctx context.Context, occasion string) ([]model.Pick, error)
	Wardrobe(ctx context.Context, loves []string, size int) ([]model.Pick, error)
}
