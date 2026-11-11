package service

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"sort"

	"github.com/victorhisamitsu/notewise/api/internal/app/perfume/domain/model"
)

// Recommender turns the current weather into a perfume suggestion.
type Recommender interface {
	Today(ctx context.Context, w model.Weather) (model.TodayPick, error)
}

type recommender struct {
	provider PerfumeProvider
	maxPicks int
	log      *slog.Logger
}

// NewRecommender creates the weather-aware recommender.
func NewRecommender(p PerfumeProvider, log *slog.Logger) Recommender {
	return &recommender{provider: p, maxPicks: 6, log: log.With("component", "recommender_service")}
}

// Conditions is the interpretation of the weather used to pick perfumes.
type Conditions struct {
	Occasion string
	Mood     string
	Headline string
	Reason   string
	Prefer   []string // family ids pushed up
	Avoid    []string // family ids pushed down
}

var (
	lightFamilies = []string{"sunlit-citrus", "green-leaf", "salt-mineral", "orchard-fruit"}
	heavyFamilies = []string{"amber-resin", "velvet-sweet", "oud-court", "smoke-leather", "spice-route"}
	warmFamilies  = []string{"amber-resin", "spice-route", "velvet-sweet", "wood-root"}
	cozyFamilies  = []string{"skin-musk", "powder-veil", "wood-root"}
)

func isRain(code int) bool {
	return (code >= 51 && code <= 67) || (code >= 80 && code <= 82) || code >= 95
}

func isSnow(code int) bool { return (code >= 71 && code <= 77) || code == 85 || code == 86 }

// Classify maps weather to an occasion and family preferences.
//
// Heat and humidity make heavy, sweet compositions cloying, while cold air
// mutes light ones — perfumers' rule of thumb, encoded as data.
func Classify(w model.Weather) Conditions {
	t := math.Round(w.TemperatureC)
	switch {
	case w.TemperatureC >= 27 && w.Humidity >= 70:
		return Conditions{"summer-days", "hot-humid", "Hot and humid: keep it light",
			fmt.Sprintf("At %.0f°C with %.0f%% humidity, heavy notes turn cloying. Citrus, green and aquatic scents stay crisp.", t, w.Humidity),
			lightFamilies, heavyFamilies}
	case w.TemperatureC >= 27:
		return Conditions{"summer-days", "hot", "A warm day calls for something bright",
			fmt.Sprintf("At %.0f°C, fresh and zesty compositions feel best on skin.", t),
			lightFamilies, heavyFamilies}
	case w.TemperatureC <= 16 || isSnow(w.Code):
		return Conditions{"cold-weather", "cold", "Cold air needs warmth",
			fmt.Sprintf("At %.0f°C, light scents fade fast. Amber, spice and woods project and linger.", t),
			warmFamilies, lightFamilies}
	case isRain(w.Code):
		return Conditions{"cold-weather", "rainy", "Rainy day, cozy scents",
			fmt.Sprintf("Rain at %.0f°C pairs with soft musks, woods and comforting sweetness.", t),
			cozyFamilies, lightFamilies}
	case !w.IsDay:
		return Conditions{"date-night", "evening", "An evening out",
			fmt.Sprintf("A mild %.0f°C evening is the moment for something deeper and more intimate.", t),
			warmFamilies, nil}
	default:
		return Conditions{"everyday", "mild", "Mild weather, easy choices",
			fmt.Sprintf("At %.0f°C almost anything works, so here are versatile everyday scents.", t),
			nil, nil}
	}
}

func (r *recommender) Today(ctx context.Context, w model.Weather) (model.TodayPick, error) {
	c := Classify(w)
	picks, err := r.provider.Edit(ctx, c.Occasion)
	if err != nil {
		r.log.ErrorContext(ctx, "edit unavailable", "occasion", c.Occasion, "err", err)
		return model.TodayPick{}, err
	}
	picks = Rank(picks, c.Prefer, c.Avoid)
	if len(picks) > r.maxPicks {
		picks = picks[:r.maxPicks]
	}
	r.log.InfoContext(ctx, "today pick", "occasion", c.Occasion, "mood", c.Mood, "temperature_c", w.TemperatureC, "humidity", w.Humidity, "picks", len(picks))
	return model.TodayPick{Occasion: c.Occasion, Mood: c.Mood, Headline: c.Headline, Reason: c.Reason, Picks: picks}, nil
}

// Rank reorders picks: preferred families first, avoided last, keeping the
// provider's order inside each group (stable sort).
func Rank(picks []model.Pick, prefer, avoid []string) []model.Pick {
	score := func(p model.Pick) int {
		for _, f := range prefer {
			if p.Family.ID == f {
				return 1
			}
		}
		for _, f := range avoid {
			if p.Family.ID == f {
				return -1
			}
		}
		return 0
	}
	out := append([]model.Pick(nil), picks...)
	sort.SliceStable(out, func(i, j int) bool { return score(out[i]) > score(out[j]) })
	return out
}
