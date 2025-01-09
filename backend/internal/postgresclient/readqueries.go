package postgresclient

import (
	"context"
	"database/sql"
	"discogsbackend/internal/models"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// FetchReleasesByYear returns a list of releases for a given year
func (pc *PostgresClient) FetchReleasesByYear(ctx context.Context, yearStr string) ([]models.DiscogsRelease, error) {
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return nil, errors.New("invalid year format")
	}

	query := `
        SELECT id, title, year, catalog_no, thumb, resource_url
        FROM releases
        WHERE year = $1
    `
	return pc.fetchReleases(ctx, query, year)
}

// FetchReleasesByArtist returns a list of releases for a given artist
func (pc *PostgresClient) FetchReleasesByArtist(ctx context.Context, artistName string) ([]models.DiscogsRelease, error) {
	return pc.fetchReleasesByRelationField(ctx, "releases_artists", "artist_id", artistName)
}

// FetchReleasesByFormat returns a list of releases for a given format
func (pc *PostgresClient) FetchReleasesByFormat(ctx context.Context, formatName string) ([]models.DiscogsRelease, error) {
	return pc.fetchReleasesByRelationField(ctx, "releases_formats", "format_id", formatName)
}

// FetchAllReleases returns a list of all releases
func (pc *PostgresClient) FetchAllReleases(ctx context.Context) ([]models.DiscogsRelease, error) {
	query := `
        SELECT r.id, r.title, r.year, r.catalog_no, r.thumb, r.resource_url
        FROM releases r
        JOIN releases_artists ra ON r.id = ra.release_id
        JOIN artists a ON a.id = ra.artist_id
    `
	return pc.fetchReleases(ctx, query)
}

// insertOrGetIDs inserts new records or retrieves existing IDs for the given values
func (pc *PostgresClient) insertOrGetIDs(ctx context.Context, tx *sql.Tx, table, column, values string, cache map[string]int) ([]int, error) {
	names := strings.Split(values, ",")
	var ids []int
	for _, name := range names {
		name = strings.TrimSpace(name)
		if id, exists := cache[name]; exists {
			ids = append(ids, id)
		} else {
			var id int
			err := tx.QueryRowContext(ctx, fmt.Sprintf(`
                INSERT INTO %s (%s)
                VALUES ($1)
                ON CONFLICT (%s) DO UPDATE SET %s = EXCLUDED.%s
                RETURNING id
            `, table, column, column, column, column), name).Scan(&id)
			if err != nil {
				return nil, fmt.Errorf("failed to insert %s '%s': %w", table, name, err)
			}
			cache[name] = id
			ids = append(ids, id)
		}
	}
	return ids, nil
}

// FetchAllArtists returns a list of all artists
func (pc *PostgresClient) FetchAllArtists(ctx context.Context) ([]models.Artist, error) {
	query := `
        SELECT id, name
        FROM artists
    `
	rows, err := pc.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var artists []models.Artist
	for rows.Next() {
		var artist models.Artist
		if err := rows.Scan(&artist.ID, &artist.Name); err != nil {
			return nil, err
		}
		artists = append(artists, artist)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return artists, nil
}

// FetchAllFormats returns a list of all formats
func (pc *PostgresClient) FetchAllFormats(ctx context.Context) ([]models.Format, error) {
	query := `
        SELECT id, name
        FROM formats
    `
	rows, err := pc.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var formats []models.Format
	for rows.Next() {
		var format models.Format
		if err := rows.Scan(&format.ID, &format.Name); err != nil {
			return nil, err
		}
		formats = append(formats, format)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return formats, nil
}

// associateIDs associates IDs in a join table
func (pc *PostgresClient) associateIDs(ctx context.Context, tx *sql.Tx, table, col1, col2 string, id1 int, ids2 []int) error {
	for _, id2 := range ids2 {
		_, err := tx.ExecContext(ctx, fmt.Sprintf(`
            INSERT INTO %s (%s, %s)
            VALUES ($1, $2)
            ON CONFLICT DO NOTHING
        `, table, col1, col2), id1, id2)
		if err != nil {
			return fmt.Errorf("failed to associate %s with %s: %w", col1, col2, err)
		}
	}
	return nil
}

// fetchReleases executes a query to fetch releases and builds the result set
func (pc *PostgresClient) fetchReleases(ctx context.Context, query string, args ...interface{}) ([]models.DiscogsRelease, error) {
	rows, err := pc.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var results []models.DiscogsRelease
	for rows.Next() {
		var release models.Release
		if err := rows.Scan(&release.ID, &release.Title, &release.Year, &release.CatalogNo, &release.Thumb, &release.ResourceURL); err != nil {
			return nil, err
		}
		results = append(results, pc.buildDiscogsRelease(ctx, release))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// fetchReleasesByRelationField fetches releases by a related field (e.g., artist or format)
func (pc *PostgresClient) fetchReleasesByRelationField(ctx context.Context, relationTable string, relationJoinColumn string, relationName string) ([]models.DiscogsRelease, error) {
	query := fmt.Sprintf(`
        SELECT r.id, r.title, r.year, r.catalog_no, r.thumb, r.resource_url
        FROM releases r
        JOIN %s rf ON r.id = rf.release_id
        JOIN %s f ON f.id = rf.%s
        WHERE f.name ILIKE $1
    `, relationTable, relationTable[:len(relationTable)-1], relationJoinColumn)

	return pc.fetchReleases(ctx, query, "%"+relationName+"%")
}

// buildDiscogsRelease constructs a DiscogsRelease from a Release
func (pc *PostgresClient) buildDiscogsRelease(ctx context.Context, release models.Release) models.DiscogsRelease {
	artists, err := pc.fetchRelatedArtists(ctx, release.ID)
	if err != nil {
		artists = []models.Artist{}
	}
	formats, err := pc.fetchRelatedFormats(ctx, release.ID)
	if err != nil {
		formats = []models.Format{}
	}

	return models.DiscogsRelease{
		ID:          release.ID,
		Title:       release.Title,
		Year:        release.Year,
		Artists:     artists,
		CatalogNo:   release.CatalogNo,
		Thumb:       release.Thumb,
		ResourceURL: release.ResourceURL,
		Formats:     formats,
		Status:      "", // Add status if needed
	}
}

// fetchRelatedArtists fetches related artists for a release
func (pc *PostgresClient) fetchRelatedArtists(ctx context.Context, releaseID int) ([]models.Artist, error) {
	query := `
        SELECT a.id, a.name
        FROM artists a
        JOIN releases_artists ra ON a.id = ra.artist_id
        WHERE ra.release_id = $1
    `

	rows, err := pc.db.QueryContext(ctx, query, releaseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var artists []models.Artist
	for rows.Next() {
		var artist models.Artist
		if err := rows.Scan(&artist.ID, &artist.Name); err != nil {
			return nil, err
		}
		artists = append(artists, artist)
	}
	return artists, nil
}

// fetchRelatedFormats fetches related formats for a release
func (pc *PostgresClient) fetchRelatedFormats(ctx context.Context, releaseID int) ([]models.Format, error) {
	query := `
        SELECT f.id, f.name
        FROM formats f
        JOIN releases_formats rf ON f.id = rf.format_id
        WHERE rf.release_id = $1
    `

	rows, err := pc.db.QueryContext(ctx, query, releaseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var formats []models.Format
	for rows.Next() {
		var format models.Format
		if err := rows.Scan(&format.ID, &format.Name); err != nil {
			return nil, err
		}
		formats = append(formats, format)
	}
	return formats, nil
}
