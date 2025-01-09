package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"discogsbackend/internal/apiclient"
	"discogsbackend/internal/postgresclient"
	"discogsbackend/internal/routes"

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
	databaseURL := viper.GetString("DATABASE_URL")

	// Initialize Discogs client
	discogsClient := apiclient.NewDiscogsClient(&http.Client{}, discogsAPIURL, discogsToken)

	// Initialize Postgres client
	ctx := context.Background()
	postgresClient, err := postgresclient.NewPostgresClient(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize Postgres client: %v", err)
	}
	defer postgresClient.Close()

	// Create database schema
	if err := postgresClient.CreateSchema(ctx); err != nil {
		log.Fatalf("Failed to create database schema: %v", err)
	}

	// Create our server with the Postgres client and Discogs client
	srv := routes.NewServer(discogsClient, postgresClient)

	// Set up the HTTP server
	httpServer := &http.Server{
		Addr:         ":8080",
		Handler:      corsMiddleware(srv.Routes()),
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	}
	srv.SetUpService("2175451")

	log.Println("Starting server on :8080")
	err = httpServer.ListenAndServe()
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
