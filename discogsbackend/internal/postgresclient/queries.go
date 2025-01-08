package postgresclient

import (
	"context"
	"database/sql"
	"discogsbackend/internal/models"
	"fmt"

	_ "github.com/lib/pq"
)

// WriteToDB inserts the fetched releases (and their related artists, styles, and genres) into the DB.
func (pc *PostgresClient) WriteToDB(ctx context.Context, releases []models.DiscogsRelease) error {
	// Start a transaction
	tx, err := pc.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ensure we roll back on error or panic
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

	for _, r := range releases {
		// Insert or get release ID
		var releaseID int
		row := tx.QueryRowContext(ctx, `
            INSERT INTO releases (title, year, label_id)
            VALUES ($1, $2, $3)
            ON CONFLICT DO NOTHING
            RETURNING id
        `, r.Title, r.Year /*someLabelID*/, 123)

		// Attempt to scan the new release ID. If the row doesn't exist (ON CONFLICT DO NOTHING), scanning might return sql.ErrNoRows.
		if scanErr := row.Scan(&releaseID); scanErr != nil {
			// If it's just sql.ErrNoRows, that means there was a conflict and no new row was inserted.
			// You might decide to select the existing release ID here if needed.
			if scanErr != sql.ErrNoRows {
				return fmt.Errorf("failed to insert release '%s': %w", r.Title, scanErr)
			}
		}

		// // Insert artists
		// for _, a := range r.Artists {
		//     // Insert or fetch artist ID
		//     // Insert bridging table (releases_artists)
		// }

		// // Insert styles
		// for _, s := range r.Styles {
		//     // Insert or fetch style ID
		//     // Insert bridging table (releases_styles)
		// }

		// // Insert genres
		// for _, g := range r.Genres {
		//     // Insert or fetch genre ID
		//     // Insert bridging table (releases_genres)
		// }
	}

	return nil
}

// FetchReleasesByArtist returns a list of releases or aggregated counts (styles) for a given artist
func (pc *PostgresClient) FetchReleasesByArtist(ctx context.Context, artistName string) ([]SomeDTO, error) {
	q := `
        SELECT st.name AS style_name,
               COUNT(DISTINCT r.id) AS release_count
          FROM releases r
          JOIN releases_artists ra ON r.id = ra.release_id
          JOIN artists a ON a.id = ra.artist_id
          JOIN releases_styles rs ON r.id = rs.release_id
          JOIN styles st ON st.id = rs.style_id
         WHERE a.name ILIKE $1
         GROUP BY st.name
    `
	rows, err := pc.db.QueryContext(ctx, q, "%"+artistName+"%")
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var results []SomeDTO
	for rows.Next() {
		var item SomeDTO
		if scanErr := rows.Scan(&item.StyleName, &item.ReleaseCount); scanErr != nil {
			return nil, scanErr
		}
		results = append(results, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// SomeDTO is a sample data transfer object for aggregated style/release info.
type SomeDTO struct {
	StyleName    string `json:"style_name"`
	ReleaseCount int    `json:"release_count"`
}
