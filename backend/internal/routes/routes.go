package routes

import (
	"discogsbackend/internal/apiclient"
	"discogsbackend/internal/postgresclient"
	"net/http"
)

type ServerInterface interface {
	Routes() http.Handler
	FetchLabelHandler(w http.ResponseWriter, r *http.Request)
	AllReleasesHandler(w http.ResponseWriter, r *http.Request)
	SearchHandler(w http.ResponseWriter, r *http.Request)
	FetchFormatsHandler(w http.ResponseWriter, r *http.Request)
	FetchArtistsHandler(w http.ResponseWriter, r *http.Request)
	CreateSchemaHandler(w http.ResponseWriter, r *http.Request)
}

// Server holds references to external dependencies
type Server struct {
	Discogs  apiclient.DiscogsClientInterface
	Postgres postgresclient.PostgresClientInterface
}

// NewServer constructs a new Server. You might pass additional dependencies here.
func NewServer(discogsClient apiclient.DiscogsClientInterface, pgClient postgresclient.PostgresClientInterface) *Server {
	return &Server{
		Discogs:  discogsClient,
		Postgres: pgClient,
	}
}

// Routes creates an http.ServeMux to handle all routes.
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	// Define route paths
	const (
		allReleasesPath  = "/all"
		searchPath       = "/search"
		fetchLabelPath   = "/fetch/label/"
		createSchemaPath = "/createschema"
		fetchFormatsPath = "/formats"
		fetchArtistsPath = "/artists"
	)

	// Register route handlers
	mux.HandleFunc(allReleasesPath, s.allReleasesHandler)
	mux.HandleFunc(searchPath, s.searchHandler)
	mux.HandleFunc(fetchLabelPath, s.fetchLabelHandler)
	mux.HandleFunc(createSchemaPath, s.createSchemaHandler)
	mux.HandleFunc(fetchFormatsPath, s.fetchFormatsHandler)
	mux.HandleFunc(fetchArtistsPath, s.fetchArtistsHandler)

	return mux
}
