// Package scento is the repository for the Scento Perfume API
// (https://www.scento.com/developers/perfume-api): a free, keyless,
// read-only perfume database.
//
// Its terms ask clients to cache answers and keep request rates
// reasonable, so every response is cached in memory (and, one level up,
// at the CDN through Cache-Control headers) and identical concurrent
// requests are coalesced into a single upstream call.
package scento

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/victorhisamitsu/notewise/api/internal/app/perfume/domain/model"
	"github.com/victorhisamitsu/notewise/api/internal/crosscutting/cache"
)

// Attribution must be shown wherever Scento data is displayed.
const Attribution = "Scento Perfume API, scento.com"

// Client calls the Scento Perfume API.
type Client struct {
	baseURL   string
	http      *http.Client
	cache     *cache.TTL[[]byte]
	userAgent string
	log       *slog.Logger
	retryWait time.Duration
}

// NewClient creates a client. cacheTTL matches Scento's own max-age (1h).
func NewClient(baseURL string, timeout, cacheTTL time.Duration, userAgent string, log *slog.Logger) *Client {
	return &Client{
		baseURL:   strings.TrimRight(baseURL, "/"),
		http:      &http.Client{Timeout: timeout},
		cache:     cache.NewTTL[[]byte](cacheTTL, 2000),
		userAgent: userAgent,
		log:       log.With("component", "scento_client"),
		retryWait: 250 * time.Millisecond,
	}
}

// Search finds perfumes by name, house, note or mood word.
func (c *Client) Search(ctx context.Context, q string) ([]model.Summary, int, error) {
	var d searchDetail
	if _, err := c.get(ctx, "/search", url.Values{"q": {q}, "kind": {"perfume"}}, &d); err != nil {
		return nil, 0, err
	}
	out := make([]model.Summary, 0, len(d.Matches))
	for _, m := range d.Matches {
		if m.Kind == "perfume" {
			out = append(out, summary(m))
		}
	}
	return out, d.Total, nil
}

// Perfume returns the full profile of one perfume.
func (c *Client) Perfume(ctx context.Context, id string) (model.Perfume, error) {
	var d perfumeDetail
	_, err := c.get(ctx, "/perfumes/"+url.PathEscape(id), nil, &d)
	if err != nil {
		return model.Perfume{}, err
	}
	if d.Found != nil && !*d.Found {
		return model.Perfume{}, model.ErrPerfumeNotFound
	}
	return toPerfume(d, c.pageURL(id)), nil
}

// Neighbours returns perfumes that share structure with id.
func (c *Client) Neighbours(ctx context.Context, id string) ([]model.Neighbour, error) {
	var d neighboursDetail
	if _, err := c.get(ctx, "/perfumes/"+url.PathEscape(id)+"/neighbours", nil, &d); err != nil {
		return nil, err
	}
	out := make([]model.Neighbour, 0, len(d.Neighbours))
	for _, n := range d.Neighbours {
		out = append(out, model.Neighbour{Summary: summary(n.perfumeRef), SharedNotes: nonNil(n.Shared), Difference: n.Difference})
	}
	return out, nil
}

// MatchNotes finds perfumes that contain all the given notes.
func (c *Client) MatchNotes(ctx context.Context, notes []string, tier string) ([]model.NoteMatch, int, error) {
	q := url.Values{"materials": {strings.Join(notes, ",")}}
	if tier != "" && tier != "any" {
		q.Set("tier", tier)
	}
	var d matchDetail
	if _, err := c.get(ctx, "/materials/match", q, &d); err != nil {
		return nil, 0, err
	}
	out := make([]model.NoteMatch, 0, len(d.Perfumes))
	for _, p := range d.Perfumes {
		tiers := p.TierHits
		if tiers == nil {
			tiers = map[string]string{}
		}
		out = append(out, model.NoteMatch{Summary: summary(p.perfumeRef), Tiers: tiers})
	}
	return out, d.Total, nil
}

// Families lists the olfactory families.
func (c *Client) Families(ctx context.Context) ([]model.Family, error) {
	var d familiesDetail
	if _, err := c.get(ctx, "/families", nil, &d); err != nil {
		return nil, err
	}
	out := make([]model.Family, 0, len(d.Tree))
	for _, f := range d.Tree {
		out = append(out, model.Family{ID: f.ID, Name: f.Name, Definition: f.Definition, Count: f.Perfumes, Subfamilies: f.Subfamilies})
	}
	return out, nil
}

// Edit returns a curated selection for an occasion (e.g. "summer-days").
func (c *Client) Edit(ctx context.Context, occasion string) ([]model.Pick, error) {
	var d editDetail
	if _, err := c.get(ctx, "/edits/"+url.PathEscape(occasion), nil, &d); err != nil {
		return nil, err
	}
	out := make([]model.Pick, 0, len(d.Edit))
	for _, e := range d.Edit {
		out = append(out, model.Pick{Summary: summary(e.perfumeRef), Why: e.Why})
	}
	return out, nil
}

// Note describes one ingredient and the perfumes that showcase it.
func (c *Client) Note(ctx context.Context, slug string) (model.Note, error) {
	var d materialDetail
	if _, err := c.get(ctx, "/materials/"+url.PathEscape(slug), nil, &d); err != nil {
		return model.Note{}, err
	}
	if d.Found != nil && !*d.Found {
		return model.Note{}, model.ErrNoteNotFound
	}
	return toNote(d), nil
}

func (c *Client) pageURL(id string) string {
	return "https://www.scento.com/perfumes/" + id
}

// get fetches path, unwraps the envelope and decodes detail into out.
func (c *Client) get(ctx context.Context, path string, q url.Values, out any) (envelope, error) {
	u := c.baseURL + path
	if len(q) > 0 {
		u += "?" + q.Encode()
	}

	body, hit, err := c.cache.GetOrLoad(ctx, u, func(ctx context.Context) ([]byte, error) {
		return c.fetch(ctx, u)
	})
	if err != nil {
		return envelope{}, err
	}
	c.log.DebugContext(ctx, "scento response", "path", path, "cache_hit", hit)

	var env envelope
	if err := json.Unmarshal(body, &env); err != nil {
		return envelope{}, fmt.Errorf("%w: decode envelope: %v", model.ErrUpstream, err)
	}
	if err := json.Unmarshal(env.Detail, out); err != nil {
		return envelope{}, fmt.Errorf("%w: decode detail: %v", model.ErrUpstream, err)
	}
	return env, nil
}

// fetch performs the HTTP call, retrying once on network errors and 5xx.
func (c *Client) fetch(ctx context.Context, u string) ([]byte, error) {
	var lastErr error
	for attempt := 0; attempt < 2; attempt++ {
		if attempt > 0 {
			select {
			case <-time.After(c.retryWait):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
		start := time.Now()
		body, status, err := c.do(ctx, u)
		elapsed := time.Since(start).Milliseconds()

		switch {
		case err != nil:
			lastErr = fmt.Errorf("%w: %v", model.ErrUpstream, err)
			c.log.WarnContext(ctx, "scento request failed", "url", u, "attempt", attempt+1, "duration_ms", elapsed, "err", err)
			continue
		case status >= 500 || status == http.StatusTooManyRequests:
			lastErr = fmt.Errorf("%w: status %d", model.ErrUpstream, status)
			c.log.WarnContext(ctx, "scento returned an error status", "url", u, "attempt", attempt+1, "status", status, "duration_ms", elapsed)
			continue
		case status == http.StatusNotFound:
			return nil, model.ErrPerfumeNotFound
		case status >= 400:
			var ae apiError
			_ = json.Unmarshal(body, &ae)
			return nil, fmt.Errorf("%w: %s", model.ErrInvalidInput, ae.Error.Message)
		}
		c.log.InfoContext(ctx, "scento request", "url", u, "status", status, "duration_ms", elapsed)
		return body, nil
	}
	return nil, lastErr
}

func (c *Client) do(ctx context.Context, u string) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if err != nil {
		return nil, resp.StatusCode, err
	}
	return body, resp.StatusCode, nil
}
