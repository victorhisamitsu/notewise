package model

// Weather is what the browser read from the weather API.
type Weather struct {
	TemperatureC float64
	Humidity     float64
	IsDay        bool
	Code         int // WMO weather code
}

// TodayPick is the weather-aware recommendation.
type TodayPick struct {
	Occasion string `json:"occasion"`
	Headline string `json:"headline"`
	Reason   string `json:"reason"`
	Mood     string `json:"mood"` // e.g. "hot-humid", "cold", "rainy"
	Picks    []Pick `json:"picks"`
}

// Count is a label with a frequency, used in profile charts.
type Count struct {
	Label string `json:"label"`
	Count int    `json:"count"`
}

// CollectionProfile summarizes the olfactory taste behind a collection.
type CollectionProfile struct {
	Size           int      `json:"size"`
	Families       []Count  `json:"families"`
	TopNotes       []Count  `json:"topNotes"`
	Seasons        []Count  `json:"seasons"`
	Hours          []Count  `json:"hours"`
	CharacterWords []Count  `json:"characterWords"`
	Headline       string   `json:"headline"`
	Gaps           []string `json:"gaps"`        // seasons / moments not covered
	Suggestions    []Pick   `json:"suggestions"` // perfumes that fill the wardrobe
	Missing        []string `json:"missing"`     // ids the provider did not know
}
