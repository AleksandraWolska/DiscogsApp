package routes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"discogsbackend/internal/models"
)

// That is pretty much duplication to below method
func (s *Server) SetUpService(labelID string) error {

	err := s.Postgres.CreateSchema(context.Background())
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to create schema: %v", err))
	}

	fetchReleases, err := s.Discogs.FetchOnLabel(labelID)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to fetch releases: %v", err))
	}

	err = s.Postgres.WriteToDB(context.Background(), fetchReleases)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to write releases to DB: %v", err))
	}

	log.Printf("Successfully fetched and stored releases for label %q\n", labelID)

	return nil
}

// This is not needed, remaining from first app version
func (s *Server) fetchLabelHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	labelID := strings.TrimPrefix(r.URL.Path, "/fetch/label/")
	if labelID == "" {
		httpError(w, "Missing label ID", http.StatusBadRequest)
		return
	}

	fetchReleases, err := s.Discogs.FetchOnLabel(labelID)
	if err != nil {
		httpError(w, fmt.Sprintf("Failed to fetch releases: %v", err), http.StatusInternalServerError)
		return
	}

	err = s.Postgres.WriteToDB(ctx, fetchReleases)
	if err != nil {
		httpError(w, fmt.Sprintf("Failed to write releases to DB: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully fetched and stored releases for label %q\n", labelID)
	fmt.Fprintf(w, "Successfully fetched and stored releases for label %q\n", labelID)
}

func (s *Server) allReleasesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		results []models.DiscogsRelease
		err     error
	)

	results, err = s.Postgres.FetchAllReleases(ctx)

	if err != nil {
		httpError(w, fmt.Sprintf("Failed to fetch releases: %v", err), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(results); err != nil {
		httpError(w, fmt.Sprintf("Failed to encode results: %v", err), http.StatusInternalServerError)
	}
}

// SQL was never my strong feature. 3 hours are way too little, this search does not work

// func (s *Server) searchHandler(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
// 	q := r.URL.Query()
// 	artist := q.Get("artist")
// 	format := q.Get("format")
// 	year := q.Get("year")

// 	var (
// 		results []models.DiscogsRelease
// 		err     error
// 	)

// 	switch {
// 	case artist != "":
// 		results, err = s.Postgres.FetchReleasesByArtist(ctx, artist)
// 	case format != "":
// 		results, err = s.Postgres.FetchReleasesByFormat(ctx, format)
// 	case year != "":
// 		results, err = s.Postgres.FetchReleasesByYear(ctx, year)
// 	default:
// 		httpError(w, "No valid query parameter provided (artist, format, or year)", http.StatusBadRequest)
// 		return
// 	}

// 	if err != nil {
// 		httpError(w, fmt.Sprintf("Failed to fetch releases: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	// Encode results as JSON
// 	if err := json.NewEncoder(w).Encode(results); err != nil {
// 		httpError(w, fmt.Sprintf("Failed to encode results: %v", err), http.StatusInternalServerError)
// 	}
// }

func (s *Server) fetchFormatsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	formats, err := s.Postgres.FetchAllFormats(ctx)
	if err != nil {
		httpError(w, fmt.Sprintf("Failed to fetch formats: %v", err), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(formats); err != nil {
		httpError(w, fmt.Sprintf("Failed to encode results: %v", err), http.StatusInternalServerError)
	}
}

func (s *Server) fetchArtistsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	artists, err := s.Postgres.FetchAllArtists(ctx)
	if err != nil {
		httpError(w, fmt.Sprintf("Failed to fetch artists: %v", err), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(artists); err != nil {
		httpError(w, fmt.Sprintf("Failed to encode results: %v", err), http.StatusInternalServerError)
	}
}

func (s *Server) createSchemaHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := s.Postgres.CreateSchema(ctx)
	if err != nil {
		httpError(w, fmt.Sprintf("Failed to create schema: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "Schema created successfully")
}

func httpError(w http.ResponseWriter, message string, code int) {
	http.Error(w, message, code)
	log.Printf("HTTP %d - %s", code, message)
}

func (s *Server) ListTablesAndRelationsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	relations, err := s.Postgres.ListTablesAndRelations(ctx)
	if err != nil {
		http.Error(w, "Failed to list tables and relations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(relations)
}
