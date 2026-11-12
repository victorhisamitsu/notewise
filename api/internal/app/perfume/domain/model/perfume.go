// Package model holds Notewise's domain entities and errors. They are
// independent of the external provider: the Scento repository maps its
// payloads into these types.
package model

import (
	"errors"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// Domain errors, translated into HTTP responses by the handler layer.
var (
	ErrPerfumeNotFound = errors.New("perfume not found")
	ErrNoteNotFound    = errors.New("note not found")
	ErrInvalidInput    = errors.New("invalid input")
	ErrUpstream        = errors.New("perfume provider unavailable")
)

// Family is one of the provider's olfactory families.
type Family struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Definition  string   `json:"definition,omitempty"`
	Count       int      `json:"count,omitempty"`
	Subfamilies []string `json:"subfamilies,omitempty"`
}

// Pyramid is the classic top / heart / base structure of a perfume.
type Pyramid struct {
	Top   []string `json:"top"`
	Heart []string `json:"heart"`
	Base  []string `json:"base"`
}

// Summary is the compact form of a perfume used in lists.
type Summary struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	House  string `json:"house"`
	Family Family `json:"family"`
}

// Perfume is the full profile of a perfume.
type Perfume struct {
	Summary
	Year           int      `json:"year,omitempty"`
	Perfumer       string   `json:"perfumer,omitempty"`
	Concentration  string   `json:"concentration,omitempty"`
	Subfamily      string   `json:"subfamily,omitempty"`
	Pyramid        *Pyramid `json:"pyramid,omitempty"` // nil when the house publishes a flat list
	Notes          []string `json:"notes"`             // every note, flattened
	Character      string   `json:"character,omitempty"`
	CharacterWords []string `json:"characterWords"`
	Seasons        []string `json:"seasons"`
	Hours          []string `json:"hours"`
	Signature      []string `json:"signatureNotes"`
	SourceURL      string   `json:"sourceUrl,omitempty"`
}

// AllNotes returns the pyramid flattened (top → base) or the flat list.
func (p Perfume) AllNotes() []string {
	if p.Pyramid == nil {
		return p.Notes
	}
	out := make([]string, 0, len(p.Pyramid.Top)+len(p.Pyramid.Heart)+len(p.Pyramid.Base))
	out = append(out, p.Pyramid.Top...)
	out = append(out, p.Pyramid.Heart...)
	return append(out, p.Pyramid.Base...)
}

// Neighbour is a perfume similar to another one.
type Neighbour struct {
	Summary
	SharedNotes []string `json:"sharedNotes"`
	Difference  string   `json:"difference,omitempty"`
}

// NoteMatch is a perfume that contains the requested notes.
type NoteMatch struct {
	Summary
	// Where each requested note sits in the pyramid: top, heart, base or listed.
	Tiers map[string]string `json:"tiers"`
}

// Note is an ingredient (raw material) as perfumers describe it.
type Note struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Kind       string   `json:"kind,omitempty"` // natural, synthetic, or both
	SmellsLike string   `json:"smellsLike,omitempty"`
	Origin     string   `json:"origin,omitempty"`
	Category   string   `json:"category,omitempty"` // e.g. Gourmand, Spice
	Family     Family   `json:"family"`
	PairsWith  []string `json:"pairsWith"`
	Aliases    []string `json:"aliases"`
	// Signature perfumes chosen by the provider, with where the note sits.
	ShowcasedIn []NoteMatch `json:"showcasedIn"`
}

// Pick is a curated suggestion with a reason.
type Pick struct {
	Summary
	Role string `json:"role,omitempty"`
	Why  string `json:"why"`
}

var (
	nonSlug   = regexp.MustCompile(`[^a-z0-9]+`)
	validSlug = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
)

// Slugify turns "Eau Sauvage" into "eau-sauvage", removing accents.
func Slugify(s string) string {
	s = strings.Map(func(r rune) rune {
		if unicode.Is(unicode.Mn, r) { // combining accent marks
			return -1
		}
		return r
	}, norm.NFD.String(strings.ToLower(s)))
	s = strings.ReplaceAll(s, "'", "")
	s = strings.ReplaceAll(s, "’", "")
	return strings.Trim(nonSlug.ReplaceAllString(s, "-"), "-")
}

// PerfumeID builds Notewise's stable id for a perfume: "<name>-<house>".
// The provider resolves this form, e.g. "eau-sauvage-dior".
func PerfumeID(name, house string) string {
	return Slugify(name + " " + house)
}

// ValidID reports whether id has the shape produced by PerfumeID.
func ValidID(id string) bool {
	return len(id) <= 120 && validSlug.MatchString(id)
}

// CleanName trims provider artefacts such as a trailing " :".
func CleanName(s string) string {
	return strings.TrimSpace(strings.TrimRight(strings.TrimSpace(s), ":"))
}
