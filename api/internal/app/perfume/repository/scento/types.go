package scento

import "encoding/json"

// envelope is the common shape of every Scento answer.
type envelope struct {
	Lede         string          `json:"lede"`
	Detail       json.RawMessage `json:"detail"`
	Cite         string          `json:"cite"`
	Continuation *string         `json:"continuation"`
}

type apiError struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type perfumeDetail struct {
	Found          *bool    `json:"found"` // present (false) only when unknown
	Name           string   `json:"name"`
	House          string   `json:"house"`
	Year           *int     `json:"year"`
	Perfumer       *string  `json:"perfumer"`
	Concentration  string   `json:"concentration"`
	Family         string   `json:"family"`
	Subfamily      *string  `json:"subfamily"`
	Pyramid        *pyramid `json:"pyramid"`
	Notes          []string `json:"notes"`
	Character      string   `json:"character"`
	CharacterWords []string `json:"character_words"`
	Wears          struct {
		Seasons []string `json:"seasons"`
		Hours   []string `json:"hours"`
	} `json:"wears"`
	SignatureMaterials []string `json:"signature_materials"`
}

type pyramid struct {
	Top   []string `json:"top"`
	Heart []string `json:"heart"`
	Base  []string `json:"base"`
}

type perfumeRef struct {
	Kind   string `json:"kind"`
	Name   string `json:"name"`
	House  string `json:"house"`
	Family string `json:"family"`
}

type searchDetail struct {
	Total   int          `json:"total"`
	Matches []perfumeRef `json:"matches"`
}

type neighboursDetail struct {
	Total      int `json:"total"`
	Neighbours []struct {
		perfumeRef
		Shared     []string `json:"shared"`
		Difference string   `json:"difference"`
	} `json:"neighbours"`
}

type matchDetail struct {
	Total    int `json:"total"`
	Perfumes []struct {
		perfumeRef
		TierHits map[string]string `json:"tier_hits"`
	} `json:"perfumes"`
}

type familiesDetail struct {
	Tree []struct {
		ID          string   `json:"id"`
		Name        string   `json:"name"`
		Definition  string   `json:"definition"`
		Subfamilies []string `json:"subfamilies"`
		Perfumes    int      `json:"perfumes"`
	} `json:"tree"`
}

type editDetail struct {
	Edit []struct {
		perfumeRef
		Why string `json:"why"`
	} `json:"edit"`
}

type wardrobeDetail struct {
	Wardrobe []struct {
		perfumeRef
		Role string `json:"role"`
		Why  string `json:"why"`
	} `json:"wardrobe"`
}
