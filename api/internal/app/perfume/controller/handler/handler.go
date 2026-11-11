// Package handler adapts HTTP requests to use cases: it parses input,
// sets cache headers and translates domain errors into status codes.
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/victorhisamitsu/notewise/api/internal/app/perfume/controller/dto"
	"github.com/victorhisamitsu/notewise/api/internal/app/perfume/controller/usecase"
	"github.com/victorhisamitsu/notewise/api/internal/app/perfume/domain/model"
	"github.com/victorhisamitsu/notewise/api/internal/crosscutting/httpx"
)

// Perfume exposes every Notewise endpoint.
type Perfume struct {
	discovery usecase.Discovery
	personal  usecase.Personal
	log       *slog.Logger
}

// NewPerfume creates the handler.
func NewPerfume(d usecase.Discovery, p usecase.Personal, log *slog.Logger) *Perfume {
	return &Perfume{discovery: d, personal: p, log: log.With("component", "perfume_handler")}
}

// Search handles GET /api/v1/perfumes/search?q=.
func (h *Perfume) Search(w http.ResponseWriter, r *http.Request) {
	res, err := h.discovery.Search(r.Context(), r.URL.Query().Get("q"))
	h.respond(w, r, res, err, httpx.CacheLong)
}

// Get handles GET /api/v1/perfumes/{id}.
func (h *Perfume) Get(w http.ResponseWriter, r *http.Request) {
	res, err := h.discovery.Perfume(r.Context(), r.PathValue("id"))
	h.respond(w, r, res, err, httpx.CacheLong)
}

// MatchNotes handles GET /api/v1/notes/match?notes=vanilla,tonka&tier=base.
func (h *Perfume) MatchNotes(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	res, err := h.discovery.MatchNotes(r.Context(), strings.Split(q.Get("notes"), ","), q.Get("tier"))
	h.respond(w, r, res, err, httpx.CacheLong)
}

// Families handles GET /api/v1/families.
func (h *Perfume) Families(w http.ResponseWriter, r *http.Request) {
	res, err := h.discovery.Families(r.Context())
	h.respond(w, r, res, err, httpx.CacheLong)
}

// Today handles GET /api/v1/recommendations/today?temp=&humidity=&isDay=&code=.
func (h *Perfume) Today(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	temp, errT := strconv.ParseFloat(q.Get("temp"), 64)
	hum, errH := strconv.ParseFloat(q.Get("humidity"), 64)
	code, errC := strconv.Atoi(q.Get("code"))
	if errT != nil || errH != nil || errC != nil {
		h.respond(w, r, nil, model.ErrInvalidInput, "")
		return
	}
	res, err := h.personal.Today(r.Context(), model.Weather{
		TemperatureC: temp, Humidity: hum, Code: code, IsDay: q.Get("isDay") != "false",
	})
	h.respond(w, r, res, err, httpx.CacheShort)
}

// Profile handles POST /api/v1/collection/profile.
func (h *Perfume) Profile(w http.ResponseWriter, r *http.Request) {
	var req dto.ProfileRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16<<10)).Decode(&req); err != nil {
		h.respond(w, r, nil, errors.Join(model.ErrInvalidInput, err), "")
		return
	}
	res, err := h.personal.Profile(r.Context(), req)
	h.respond(w, r, res, err, httpx.NoStore)
}

func (h *Perfume) respond(w http.ResponseWriter, r *http.Request, body any, err error, cache string) {
	if err == nil {
		w.Header().Set("Cache-Control", cache)
		httpx.WriteJSON(w, http.StatusOK, body)
		return
	}
	w.Header().Set("Cache-Control", httpx.NoStore)
	status, code := Status(err)
	msg := err.Error()
	if status >= 500 {
		h.log.ErrorContext(r.Context(), "request failed", "path", r.URL.Path, "status", status, "err", err)
		msg = "the perfume database is unavailable right now, please try again"
	}
	httpx.WriteJSON(w, status, dto.ErrorResponse{Error: code, Message: msg})
}

// Status maps domain errors to HTTP status codes.
func Status(err error) (int, string) {
	switch {
	case errors.Is(err, model.ErrInvalidInput):
		return http.StatusBadRequest, "invalid_input"
	case errors.Is(err, model.ErrPerfumeNotFound):
		return http.StatusNotFound, "not_found"
	case errors.Is(err, context.DeadlineExceeded):
		return http.StatusGatewayTimeout, "timeout"
	case errors.Is(err, model.ErrUpstream):
		return http.StatusBadGateway, "upstream_unavailable"
	default:
		return http.StatusInternalServerError, "internal_error"
	}
}
