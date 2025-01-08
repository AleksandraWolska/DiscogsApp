package postgresclient

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

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
