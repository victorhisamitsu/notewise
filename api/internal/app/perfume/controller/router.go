// Package controller exposes Notewise over HTTP.
package controller

import (
	"net/http"

	"github.com/victorhisamitsu/notewise/api/internal/app/perfume/controller/handler"
	"github.com/victorhisamitsu/notewise/api/internal/crosscutting/httpx"
)

// Route is one HTTP endpoint. Keeping routes in a table lets a test check
// them against the OpenAPI contract (api/openapi.yaml).
type Route struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

// Routes lists every public endpoint.
func Routes(h *handler.Perfume) []Route {
	return []Route{
		{http.MethodGet, "/api/health", health},
		{http.MethodGet, "/api/v1/perfumes/search", h.Search},
		{http.MethodGet, "/api/v1/perfumes/{id}", h.Get},
		{http.MethodGet, "/api/v1/notes/match", h.MatchNotes},
		{http.MethodGet, "/api/v1/notes/{note}", h.Note},
		{http.MethodGet, "/api/v1/families", h.Families},
		{http.MethodGet, "/api/v1/recommendations/today", h.Today},
	}
}

// RegisterRoutes mounts every route on mux.
func RegisterRoutes(mux *http.ServeMux, h *handler.Perfume) {
	for _, r := range Routes(h) {
		mux.HandleFunc(r.Method+" "+r.Path, r.Handler)
	}
}

func health(w http.ResponseWriter, _ *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
