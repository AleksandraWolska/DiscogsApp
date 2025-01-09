package postgresclient

import (
	"context"
	"discogsbackend/internal/models"
	"fmt"
)

func (pc *PostgresClient) CreateSchema(ctx context.Context) error {
	schema := `
    CREATE TABLE IF NOT EXISTS artists (
        id SERIAL PRIMARY KEY,
        name TEXT UNIQUE NOT NULL
    );

    CREATE TABLE IF NOT EXISTS formats (
        id SERIAL PRIMARY KEY,
        name TEXT UNIQUE NOT NULL
    );

    CREATE TABLE IF NOT EXISTS releases (
        id INT PRIMARY KEY,
        title TEXT NOT NULL,
        year INT,
        catalog_no TEXT,
        thumb TEXT,
        resource_url TEXT
    );

    CREATE TABLE IF NOT EXISTS releases_artists (
        release_id INT NOT NULL,
        artist_id INT NOT NULL,
        PRIMARY KEY (release_id, artist_id),
        FOREIGN KEY (release_id) REFERENCES releases (id) ON DELETE CASCADE,
        FOREIGN KEY (artist_id) REFERENCES artists (id) ON DELETE CASCADE
    );

    CREATE TABLE IF NOT EXISTS releases_formats (
        release_id INT NOT NULL,
        format_id INT NOT NULL,
        PRIMARY KEY (release_id, format_id),
        FOREIGN KEY (release_id) REFERENCES releases (id) ON DELETE CASCADE,
        FOREIGN KEY (format_id) REFERENCES formats (id) ON DELETE CASCADE
    );
    `

	_, err := pc.db.ExecContext(ctx, schema)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	return nil
}

func (pc *PostgresClient) WriteToDB(ctx context.Context, releases []models.DiscogsApiRelease) error {
	tx, err := pc.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	artistNameToID := make(map[string]int)
	formatNameToID := make(map[string]int)

	for _, r := range releases {
		artistIDs, err := pc.insertOrGetIDs(ctx, tx, "artists", "name", r.Artist, artistNameToID)
		if err != nil {
			return err
		}

		formatIDs, err := pc.insertOrGetIDs(ctx, tx, "formats", "name", r.Format, formatNameToID)
		if err != nil {
			return err
		}

		var releaseID int
		err = tx.QueryRowContext(ctx, `
            INSERT INTO releases (id, title, year, catalog_no, thumb, resource_url)
            VALUES ($1, $2, $3, $4, $5, $6)
            ON CONFLICT (id) DO UPDATE SET
                title = EXCLUDED.title,
                year = EXCLUDED.year,
                catalog_no = EXCLUDED.catalog_no,
                thumb = EXCLUDED.thumb,
                resource_url = EXCLUDED.resource_url
            RETURNING id
        `, r.ID, r.Title, r.Year, r.CatalogNo, r.ThumbnailURL, r.ResourceURL).Scan(&releaseID)
		if err != nil {
			return fmt.Errorf("failed to insert release '%s': %w", r.Title, err)
		}

		if err := pc.associateIDs(ctx, tx, "releases_artists", "release_id", "artist_id", releaseID, artistIDs); err != nil {
			return err
		}

		if err := pc.associateIDs(ctx, tx, "releases_formats", "release_id", "format_id", releaseID, formatIDs); err != nil {
			return err
		}
	}

	return nil
}
