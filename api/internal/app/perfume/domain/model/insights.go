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
