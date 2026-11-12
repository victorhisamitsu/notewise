// Package controller exposes Notewise over HTTP.
package controller

import (
	"net/http"

	"github.com/victorhisamitsu/notewise/api/internal/app/perfume/controller/handler"
)

// RegisterRoutes mounts the v1 API.
func RegisterRoutes(mux *http.ServeMux, h *handler.Perfume) {
	mux.HandleFunc("GET /api/v1/perfumes/search", h.Search)
	mux.HandleFunc("GET /api/v1/perfumes/{id}", h.Get)
	mux.HandleFunc("GET /api/v1/notes/match", h.MatchNotes)
	mux.HandleFunc("GET /api/v1/notes/{note}", h.Note)
	mux.HandleFunc("GET /api/v1/families", h.Families)
	mux.HandleFunc("GET /api/v1/recommendations/today", h.Today)
}
