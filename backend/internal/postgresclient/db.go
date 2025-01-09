package postgresclient

import (
	"context"
	"database/sql"
	"discogsbackend/internal/models"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type PostgresClientInterface interface {
	CreateSchema(ctx context.Context) error
	WriteToDB(ctx context.Context, releases []models.DiscogsApiRelease) error
	FetchAllReleases(ctx context.Context) ([]models.DiscogsRelease, error)
	FetchReleasesByArtist(ctx context.Context, artistName string) ([]models.DiscogsRelease, error)
	FetchReleasesByFormat(ctx context.Context, formatName string) ([]models.DiscogsRelease, error)
	FetchReleasesByYear(ctx context.Context, yearStr string) ([]models.DiscogsRelease, error)
	FetchAllArtists(ctx context.Context) ([]models.Artist, error)
	FetchAllFormats(ctx context.Context) ([]models.Format, error)
	Close() error
}

// PostgresClient holds a reference to the sql.DB object
type PostgresClient struct {
	db *sql.DB
}

// NewPostgresClient creates a new PostgresClient using database/sql
func NewPostgresClient(ctx context.Context, connStr string) (*PostgresClient, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	// Test the connection
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	log.Println("DB Connection established successfully")

	return &PostgresClient{
		db: db,
	}, nil
}

// Close terminates the connection
func (pc *PostgresClient) Close() error {
	return pc.db.Close()
}
