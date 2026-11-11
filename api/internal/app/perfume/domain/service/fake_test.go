package service

import (
	"context"
	"io"
	"log/slog"
	"sync"

	"github.com/victorhisamitsu/notewise/api/internal/app/perfume/domain/model"
)

var discard = slog.New(slog.NewTextHandler(io.Discard, nil))

type fakeProvider struct {
	mu            sync.Mutex
	perfumes      map[string]model.Perfume
	edit          []model.Pick
	wardrobe      []model.Pick
	neighboursErr error
	perfumeErr    error
	editOccasion  string
	inFlight      int
	maxInFlight   int
}

func (f *fakeProvider) Search(context.Context, string) ([]model.Summary, int, error) {
	return nil, 0, nil
}

func (f *fakeProvider) Perfume(_ context.Context, id string) (model.Perfume, error) {
	f.mu.Lock()
	f.inFlight++
	if f.inFlight > f.maxInFlight {
		f.maxInFlight = f.inFlight
	}
	f.mu.Unlock()
	defer func() { f.mu.Lock(); f.inFlight--; f.mu.Unlock() }()

	if f.perfumeErr != nil {
		return model.Perfume{}, f.perfumeErr
	}
	p, ok := f.perfumes[id]
	if !ok {
		return model.Perfume{}, model.ErrPerfumeNotFound
	}
	return p, nil
}

func (f *fakeProvider) Neighbours(context.Context, string) ([]model.Neighbour, error) {
	if f.neighboursErr != nil {
		return nil, f.neighboursErr
	}
	return []model.Neighbour{{Summary: model.Summary{ID: "n1"}}}, nil
}

func (f *fakeProvider) MatchNotes(context.Context, []string, string) ([]model.NoteMatch, int, error) {
	return nil, 0, nil
}

func (f *fakeProvider) Families(context.Context) ([]model.Family, error) { return nil, nil }

func (f *fakeProvider) Edit(_ context.Context, occasion string) ([]model.Pick, error) {
	f.editOccasion = occasion
	return f.edit, nil
}

func (f *fakeProvider) Wardrobe(context.Context, []string, int) ([]model.Pick, error) {
	return f.wardrobe, nil
}

func perfume(id, family string, notes []string, seasons, hours []string) model.Perfume {
	return model.Perfume{
		Summary: model.Summary{ID: id, Name: id, Family: model.Family{ID: model.Slugify(family), Name: family}},
		Notes:   notes, Seasons: seasons, Hours: hours,
	}
}

func pick(id, familyID string) model.Pick {
	return model.Pick{Summary: model.Summary{ID: id, Family: model.Family{ID: familyID}}}
}
