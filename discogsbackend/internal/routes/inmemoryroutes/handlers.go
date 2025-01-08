package inmemoryroutes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"discogsbackend/internal/apiclient"
	"discogsbackend/internal/inmemorydatabase"
)

// Server holds references to external dependencies
type Server struct {
	InMemory *inmemorydatabase.InMemoryClient
	Discogs  *apiclient.DiscogsClient
}

// NewServer constructs a new Server. You might pass additional dependencies here.
func NewServer(imClient *inmemorydatabase.InMemoryClient, discogsClient *apiclient.DiscogsClient) *Server {
	return &Server{
		InMemory: imClient,
		Discogs:  discogsClient,
	}
}

// Routes creates an http.ServeMux to handle all routes.
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	// /search?artist=...
	mux.HandleFunc("/search", s.searchHandler)

	// /fetch/label/123
	mux.HandleFunc("/fetch/label/", s.fetchLabelHandler)

	// /dbstate
	mux.HandleFunc("/dbstate", s.dbStateHandler)

	return mux
}

// fetchLabelHandler fetches a label ID from the URL path, calls the Discogs API,
// then writes the fetched releases to the in-memory database.
func (s *Server) fetchLabelHandler(w http.ResponseWriter, r *http.Request) {
	// Our route is "/fetch/label/{labelID}". We can parse {labelID} from the path.
	// Example path: "/fetch/label/123"
	path := strings.TrimPrefix(r.URL.Path, "/fetch/label/")
	if path == "" {
		http.Error(w, "Missing label ID", http.StatusBadRequest)
		return
	}

	labelID := path // or convert to int if needed

	// Call the Discogs API
	fetchReleases, err := s.Discogs.FetchOnLabel(labelID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the fetched releases to the in-memory database
	err = s.InMemory.WriteToDB(fetchReleases)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	fmt.Fprintf(w, "Successfully fetched and stored releases for label %q\n", labelID)
}

// searchHandler handles queries by artist, format, title, or year, returning aggregated results.
func (s *Server) searchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	artist := q.Get("artist")
	format := q.Get("format")
	year := q.Get("year")

	// For demonstration, if searching by artist:
	if artist != "" {
		results, err := s.InMemory.FetchReleasesByArtist(artist)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Encode results as JSON
		json.NewEncoder(w).Encode(results)
		return
	}

	// Handle search by format
	if format != "" {
		results, err := s.InMemory.FetchReleasesByFormat(format)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Encode results as JSON
		json.NewEncoder(w).Encode(results)
		return
	}

	// Handle search by year
	if year != "" {
		results, err := s.InMemory.FetchReleasesByYear(year)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Encode results as JSON
		json.NewEncoder(w).Encode(results)
		return
	}

	http.Error(w, "No valid query parameter provided (artist, format, title, or year)", http.StatusBadRequest)
}

// dbStateHandler returns the entire state of the in-memory database
func (s *Server) dbStateHandler(w http.ResponseWriter, r *http.Request) {
	state := s.InMemory.GetState()
	json.NewEncoder(w).Encode(state)
}

// For illustration, you could optionally add a helper function to parse integer IDs from the URL path:
func parseID(strID string) (int, error) {
	var id int
	_, err := fmt.Sscanf(strID, "%d", &id)
	if err != nil {
		return 0, fmt.Errorf("invalid ID format: %w", err)
	}
	return id, nil
}
