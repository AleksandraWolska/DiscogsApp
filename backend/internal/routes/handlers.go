package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"discogsbackend/internal/models"
)

// fetchLabelHandler fetches a label ID from the URL path, calls the Discogs API,
// then writes the fetched releases to the PostgreSQL database.
func (s *Server) fetchLabelHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	labelID := strings.TrimPrefix(r.URL.Path, "/fetch/label/")
	if labelID == "" {
		httpError(w, "Missing label ID", http.StatusBadRequest)
		return
	}

	// Call the Discogs API
	fetchReleases, err := s.Discogs.FetchOnLabel(labelID)
	if err != nil {
		httpError(w, fmt.Sprintf("Failed to fetch releases: %v", err), http.StatusInternalServerError)
		return
	}

	// Write the fetched releases to the PostgreSQL database
	err = s.Postgres.WriteToDB(ctx, fetchReleases)
	if err != nil {
		httpError(w, fmt.Sprintf("Failed to write releases to DB: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	log.Printf("Successfully fetched and stored releases for label %q\n", labelID)
	fmt.Fprintf(w, "Successfully fetched and stored releases for label %q\n", labelID)
}

// allReleasesHandler handles fetching all releases
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

	// Encode results as JSON
	if err := json.NewEncoder(w).Encode(results); err != nil {
		httpError(w, fmt.Sprintf("Failed to encode results: %v", err), http.StatusInternalServerError)
	}
}

// searchHandler handles queries by artist, format, or year, returning aggregated results.
func (s *Server) searchHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := r.URL.Query()
	artist := q.Get("artist")
	format := q.Get("format")
	year := q.Get("year")

	var (
		results []models.DiscogsRelease
		err     error
	)

	switch {
	case artist != "":
		results, err = s.Postgres.FetchReleasesByArtist(ctx, artist)
	case format != "":
		results, err = s.Postgres.FetchReleasesByFormat(ctx, format)
	case year != "":
		results, err = s.Postgres.FetchReleasesByYear(ctx, year)
	default:
		httpError(w, "No valid query parameter provided (artist, format, or year)", http.StatusBadRequest)
		return
	}

	if err != nil {
		httpError(w, fmt.Sprintf("Failed to fetch releases: %v", err), http.StatusInternalServerError)
		return
	}

	// Encode results as JSON
	if err := json.NewEncoder(w).Encode(results); err != nil {
		httpError(w, fmt.Sprintf("Failed to encode results: %v", err), http.StatusInternalServerError)
	}
}

// fetchFormatsHandler handles fetching all formats
func (s *Server) fetchFormatsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	formats, err := s.Postgres.FetchAllFormats(ctx)
	if err != nil {
		httpError(w, fmt.Sprintf("Failed to fetch formats: %v", err), http.StatusInternalServerError)
		return
	}

	// Encode results as JSON
	if err := json.NewEncoder(w).Encode(formats); err != nil {
		httpError(w, fmt.Sprintf("Failed to encode results: %v", err), http.StatusInternalServerError)
	}
}

// fetchArtistsHandler handles fetching all artists
func (s *Server) fetchArtistsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	artists, err := s.Postgres.FetchAllArtists(ctx)
	if err != nil {
		httpError(w, fmt.Sprintf("Failed to fetch artists: %v", err), http.StatusInternalServerError)
		return
	}

	// Encode results as JSON
	if err := json.NewEncoder(w).Encode(artists); err != nil {
		httpError(w, fmt.Sprintf("Failed to encode results: %v", err), http.StatusInternalServerError)
	}
}

// createSchemaHandler creates the necessary tables and schemas in the PostgreSQL database.
func (s *Server) createSchemaHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := s.Postgres.CreateSchema(ctx)
	if err != nil {
		httpError(w, fmt.Sprintf("Failed to create schema: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "Schema created successfully")
}

// httpError sends an HTTP error response with the specified message and status code.
func httpError(w http.ResponseWriter, message string, code int) {
	http.Error(w, message, code)
	log.Printf("HTTP %d - %s", code, message)
}

// // parseID parses an integer ID from the URL path.
// func parseID(strID string) (int, error) {
// 	var id int
// 	_, err := fmt.Sscanf(strID, "%d", &id)
// 	if err != nil {
// 		return 0, fmt.Errorf("invalid ID format: %w", err)
// 	}
// 	return id, nil
// }
