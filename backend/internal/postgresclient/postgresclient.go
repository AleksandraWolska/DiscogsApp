package postgresclient

import (
	"context"
	"database/sql"
	"discogsbackend/internal/models"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// We have relation on release_artist and release_format
// But SQL was mostly generated, sorry, I didn't work with SQL for a while

type PostgresClientInterface interface {
	CreateSchema(ctx context.Context) error
	WriteToDB(ctx context.Context, releases []models.DiscogsApiRelease) error
	FetchAllReleases(ctx context.Context) ([]models.DiscogsRelease, error)
	// FetchReleasesByArtist(ctx context.Context, artistName string) ([]models.DiscogsRelease, error)
	// FetchReleasesByFormat(ctx context.Context, formatName string) ([]models.DiscogsRelease, error)
	// FetchReleasesByYear(ctx context.Context, yearStr string) ([]models.DiscogsRelease, error)
	FetchAllArtists(ctx context.Context) ([]models.Artist, error)
	FetchAllFormats(ctx context.Context) ([]models.Format, error)
	ListTablesAndRelations(ctx context.Context) (map[string][]string, error)
	Close() error
}

type PostgresClient struct {
	db *sql.DB
}

func NewPostgresClient(ctx context.Context, connStr string) (*PostgresClient, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	log.Println("DB Connection established successfully")

	return &PostgresClient{
		db: db,
	}, nil
}

func (pc *PostgresClient) Close() error {
	return pc.db.Close()
}
