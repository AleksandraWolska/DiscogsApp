package main

import (
	"log"
	"net/http"
	"time"

	"discogsbackend/internal/apiclient"
	"discogsbackend/internal/inmemorydatabase"
	"discogsbackend/internal/routes/inmemoryroutes"

	"github.com/spf13/viper"
)

func main() {
	// Initialize Viper
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// Read configuration values
	discogsAPIURL := viper.GetString("DISCOGS_API_URL")
	discogsToken := viper.GetString("DISCOGS_TOKEN")

	// Initialize InMemory client
	imClient := inmemorydatabase.NewInMemoryClient()

	// Initialize Discogs client
	discogsClient := apiclient.NewDiscogsClient(&http.Client{}, discogsAPIURL, discogsToken)

	// Create our server with the InMemory client and Discogs client
	srv := inmemoryroutes.NewServer(imClient, discogsClient)

	// Set up the HTTP server
	httpServer := &http.Server{
		Addr:         ":8080",
		Handler:      corsMiddleware(srv.Routes()), // Our net/http mux with CORS middleware
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	}

	log.Println("Starting server on :8080")
	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// corsMiddleware adds CORS headers to the response
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
	})
}
