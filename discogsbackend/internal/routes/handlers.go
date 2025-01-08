package routes

// // Server holds references to external dependencies
// type Server struct {
// 	Postgres *postgresclient.PostgresClient
// 	Discogs  *apiclient.DiscogsClient
// }

// // NewServer constructs a new Server. You might pass additional dependencies here.
// func NewServer(pgClient *postgresclient.PostgresClient, discogsClient *apiclient.DiscogsClient) *Server {
// 	return &Server{
// 		Postgres: pgClient,
// 		Discogs:  discogsClient,
// 	}
// }

// // Routes creates an http.ServeMux to handle all routes.
// func (s *Server) Routes() http.Handler {
// 	mux := http.NewServeMux()

// 	// /search?artist=...
// 	mux.HandleFunc("/search", s.searchHandler)

// 	// /fetch/label/123
// 	mux.HandleFunc("/fetch/label/", s.fetchLabelHandler)

// 	return mux
// }

// // fetchLabelHandler fetches a label ID from the URL path, calls the Discogs API,
// // then writes the fetched releases to the DB.
// func (s *Server) fetchLabelHandler(w http.ResponseWriter, r *http.Request) {
// 	// Our route is "/fetch/label/{labelID}". We can parse {labelID} from the path.
// 	// Example path: "/fetch/label/123"
// 	path := strings.TrimPrefix(r.URL.Path, "/fetch/label/")
// 	if path == "" {
// 		http.Error(w, "Missing label ID", http.StatusBadRequest)
// 		return
// 	}

// 	labelID := path // or convert to int if needed

// 	// Call the Discogs API
// 	fetchReleases, err := s.Discogs.FetchOnLabel(labelID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Write the fetched releases to the DB
// 	err = s.Postgres.WriteToDB(r.Context(), fetchReleases)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Respond with a success message
// 	fmt.Fprintf(w, "Successfully fetched and stored releases for label %q\n", labelID)
// }

// // searchHandler handles queries by artist, style, or genre, returning aggregated results.
// func (s *Server) searchHandler(w http.ResponseWriter, r *http.Request) {
// 	q := r.URL.Query()
// 	artist := q.Get("artist")
// 	style := q.Get("style")
// 	genre := q.Get("genre")

// 	// For demonstration, if searching by artist:
// 	if artist != "" {
// 		results, err := s.Postgres.FetchReleasesByArtist(r.Context(), artist)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		// Encode results as JSON
// 		json.NewEncoder(w).Encode(results)
// 		return
// 	}

// 	// Similarly handle style or genre...
// 	if style != "" {
// 		// Example: s.Postgres.FetchReleasesByStyle(r.Context(), style)
// 		// ...
// 		fmt.Fprintf(w, "Search by style=%q not yet implemented\n", style)
// 		return
// 	}

// 	if genre != "" {
// 		// Example: s.Postgres.FetchReleasesByGenre(r.Context(), genre)
// 		// ...
// 		fmt.Fprintf(w, "Search by genre=%q not yet implemented\n", genre)
// 		return
// 	}

// 	http.Error(w, "No valid query parameter provided (artist, style, or genre)", http.StatusBadRequest)
// }

// // For illustration, you could optionally add a helper function to parse integer IDs from the URL path:
// func parseID(strID string) (int, error) {
// 	var id int
// 	_, err := fmt.Sscanf(strID, "%d", &id)
// 	if err != nil {
// 		return 0, fmt.Errorf("invalid ID format: %w", err)
// 	}
// 	return id, nil
// }
