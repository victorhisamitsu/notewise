package scento

import (
	"strings"

	"github.com/victorhisamitsu/notewise/api/internal/app/perfume/domain/model"
)

func family(name string) model.Family {
	return model.Family{ID: model.Slugify(name), Name: name}
}

func summary(r perfumeRef) model.Summary {
	name := model.CleanName(r.Name)
	return model.Summary{
		ID:     model.PerfumeID(name, r.House),
		Name:   name,
		House:  r.House,
		Family: family(r.Family),
	}
}

func toPerfume(d perfumeDetail, sourceURL string) model.Perfume {
	p := model.Perfume{
		Summary:        summary(perfumeRef{Name: d.Name, House: d.House, Family: d.Family}),
		Concentration:  d.Concentration,
		Character:      d.Character,
		CharacterWords: nonNil(d.CharacterWords),
		Seasons:        nonNil(d.Wears.Seasons),
		Hours:          nonNil(d.Wears.Hours),
		Signature:      nonNil(d.SignatureMaterials),
		SourceURL:      sourceURL,
	}
	if d.Year != nil {
		p.Year = *d.Year
	}
	if d.Perfumer != nil {
		p.Perfumer = *d.Perfumer
	}
	if d.Subfamily != nil {
		p.Subfamily = *d.Subfamily
	}
	if d.Pyramid != nil && (len(d.Pyramid.Top)+len(d.Pyramid.Heart)+len(d.Pyramid.Base)) > 0 {
		p.Pyramid = &model.Pyramid{
			Top:   nonNil(d.Pyramid.Top),
			Heart: nonNil(d.Pyramid.Heart),
			Base:  nonNil(d.Pyramid.Base),
		}
	}
	p.Notes = nonNil(p.AllNotes())
	if p.Pyramid == nil {
		p.Notes = nonNil(d.Notes)
	}
	return p
}

func str(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func toNote(d materialDetail) model.Note {
	n := model.Note{
		ID:         model.Slugify(d.Name),
		Name:       d.Name,
		Kind:       str(d.Kind),
		SmellsLike: str(d.SmellsLike),
		Origin:     str(d.Origin),
		Category:   str(d.Category),
		PairsWith:  nonNil(d.PairsWith),
		Aliases:    nonNil(d.Aliases),
	}
	if fam := str(d.ScentoFamily); fam != "" {
		n.Family = family(fam)
	}
	if len(n.Aliases) > 8 {
		n.Aliases = n.Aliases[:8]
	}
	n.ShowcasedIn = make([]model.NoteMatch, 0, len(d.ShowcasedIn))
	for _, s := range d.ShowcasedIn {
		n.ShowcasedIn = append(n.ShowcasedIn, model.NoteMatch{
			Summary: summary(perfumeRef{Name: s.Name, House: s.House}),
			Tiers:   map[string]string{d.Name: s.Tier},
		})
	}
	return n
}

func nonNil(s []string) []string {
	if s == nil {
		return []string{}
	}
	out := make([]string, 0, len(s))
	for _, v := range s {
		if v = strings.TrimSpace(v); v != "" {
			out = append(out, v)
		}
	}
	return out
}
