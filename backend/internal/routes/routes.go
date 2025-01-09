package routes

import (
	"discogsbackend/internal/apiclient"
	"discogsbackend/internal/postgresclient"
	"net/http"
)

// This should be 'handler' interface....
type ServerInterface interface {
	Routes() http.Handler
	FetchLabelHandler(w http.ResponseWriter, r *http.Request)
	AllReleasesHandler(w http.ResponseWriter, r *http.Request)
	SearchHandler(w http.ResponseWriter, r *http.Request)
	FetchFormatsHandler(w http.ResponseWriter, r *http.Request)
	FetchArtistsHandler(w http.ResponseWriter, r *http.Request)
	CreateSchemaHandler(w http.ResponseWriter, r *http.Request)
}

// ... I am pretty sure this should not be in routes package, and handler should not hold that, its more for a service struct
type Server struct {
	Discogs  apiclient.DiscogsClientInterface
	Postgres postgresclient.PostgresClientInterface
}

func NewServer(discogsClient apiclient.DiscogsClientInterface, pgClient postgresclient.PostgresClientInterface) *Server {
	return &Server{
		Discogs:  discogsClient,
		Postgres: pgClient,
	}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	const (
		allReleasesPath          = "/all"
		searchPath               = "/search"
		fetchLabelPath           = "/fetch/label/"
		createSchemaPath         = "/createschema"
		fetchFormatsPath         = "/formats"
		fetchArtistsPath         = "/artists"
		fetchRelationsSchemaPath = "/relations"
	)

	// Register route handlers
	mux.HandleFunc(allReleasesPath, s.allReleasesHandler)
	// mux.HandleFunc(searchPath, s.searchHandler)
	mux.HandleFunc(fetchLabelPath, s.fetchLabelHandler)
	mux.HandleFunc(createSchemaPath, s.createSchemaHandler)
	mux.HandleFunc(fetchFormatsPath, s.fetchFormatsHandler)
	mux.HandleFunc(fetchArtistsPath, s.fetchArtistsHandler)
	mux.HandleFunc(fetchRelationsSchemaPath, s.ListTablesAndRelationsHandler)

	return mux
}
