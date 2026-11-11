// Package dto defines Notewise's public API contracts.
package dto

import "github.com/victorhisamitsu/notewise/api/internal/app/perfume/domain/model"

// SearchResponse is returned by GET /api/v1/perfumes/search.
type SearchResponse struct {
	Query       string          `json:"query"`
	Total       int             `json:"total"`
	Results     []model.Summary `json:"results"`
	Attribution string          `json:"attribution"`
}

// PerfumeResponse is returned by GET /api/v1/perfumes/{id}.
type PerfumeResponse struct {
	Perfume     model.Perfume     `json:"perfume"`
	Neighbours  []model.Neighbour `json:"neighbours"`
	Attribution string            `json:"attribution"`
}

// NoteMatchResponse is returned by GET /api/v1/notes/match.
type NoteMatchResponse struct {
	Notes       []string          `json:"notes"`
	Tier        string            `json:"tier"`
	Total       int               `json:"total"`
	Results     []model.NoteMatch `json:"results"`
	Attribution string            `json:"attribution"`
}

// FamiliesResponse is returned by GET /api/v1/families.
type FamiliesResponse struct {
	Families    []model.Family `json:"families"`
	Attribution string         `json:"attribution"`
}

// TodayResponse is returned by GET /api/v1/recommendations/today.
type TodayResponse struct {
	model.TodayPick
	Attribution string `json:"attribution"`
}

// ProfileRequest is the body of POST /api/v1/collection/profile.
type ProfileRequest struct {
	PerfumeIDs []string `json:"perfumeIds"`
}

// ProfileResponse is returned by POST /api/v1/collection/profile.
type ProfileResponse struct {
	model.CollectionProfile
	Attribution string `json:"attribution"`
}

// ErrorResponse is the body of every error.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
