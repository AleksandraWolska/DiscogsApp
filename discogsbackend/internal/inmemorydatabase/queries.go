package inmemorydatabase

import (
	"discogsbackend/internal/models"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func (imc *InMemoryClient) WriteToDB(releases []models.DiscogsApiRelease) error {
	imc.mu.Lock()
	defer imc.mu.Unlock()

	artistNameToID := make(map[string]int)
	formatNameToID := make(map[string]int)

	for _, r := range releases {
		// Parse comma-separated artists
		artistNames := strings.Split(r.Artist, ",")
		var artistIDs []int
		for _, name := range artistNames {
			name = strings.TrimSpace(name)
			if id, exists := artistNameToID[name]; exists {
				artistIDs = append(artistIDs, id)
			} else {
				artist := models.Artist{
					ID:   imc.nextID,
					Name: name,
				}
				imc.artists[imc.nextID] = artist
				artistNameToID[name] = imc.nextID
				artistIDs = append(artistIDs, imc.nextID)
				imc.nextID++
			}
		}

		// Parse comma-separated formats
		formatNames := strings.Split(r.Format, ",")
		var formatIDs []int
		for _, name := range formatNames {
			name = strings.TrimSpace(name)
			if id, exists := formatNameToID[name]; exists {
				formatIDs = append(formatIDs, id)
			} else {
				format := models.Format{
					ID:   imc.nextID,
					Name: name,
				}
				imc.formats[imc.nextID] = format
				formatNameToID[name] = imc.nextID
				formatIDs = append(formatIDs, imc.nextID)
				imc.nextID++
			}
		}

		// Check for duplicate release
		if _, exists := imc.releases[r.ID]; exists {
			continue // Skip duplicate release
		}

		// Save release
		release := models.Release{
			ID:          r.ID,
			Title:       r.Title,
			Year:        r.Year,
			CatalogNo:   r.CatalogNo,
			Thumb:       r.ThumbnailURL,
			ResourceURL: r.ResourceURL,
			// Add other fields as needed
		}
		imc.releases[r.ID] = release

		// Associate artists with the release
		for _, artistID := range artistIDs {
			imc.releasesArtists[release.ID] = append(imc.releasesArtists[release.ID], artistID)
		}

		// Associate formats with the release
		for _, formatID := range formatIDs {
			imc.releasesFormats[release.ID] = append(imc.releasesFormats[release.ID], formatID)
		}
	}
	imc.printState()
	return nil
}

func (imc *InMemoryClient) FetchReleasesByArtist(artistName string) ([]models.DiscogsRelease, error) {
	imc.mu.Lock()
	defer imc.mu.Unlock()

	var results []models.DiscogsRelease
	for releaseID, artistIDs := range imc.releasesArtists {
		for _, artistID := range artistIDs {
			if artist, exists := imc.artists[artistID]; exists && artist.Name == artistName {
				if release, exists := imc.releases[releaseID]; exists {
					results = append(results, imc.buildDiscogsRelease(release))
				}
			}
		}
	}
	if len(results) == 0 {
		return nil, errors.New("no releases found for the specified artist")
	}
	imc.printState()
	return results, nil
}

func (imc *InMemoryClient) FetchReleasesByFormat(formatName string) ([]models.DiscogsRelease, error) {
	imc.mu.Lock()
	defer imc.mu.Unlock()

	var results []models.DiscogsRelease
	for releaseID, formatIDs := range imc.releasesFormats {
		for _, formatID := range formatIDs {
			if format, exists := imc.formats[formatID]; exists && format.Name == formatName {
				if release, exists := imc.releases[releaseID]; exists {
					results = append(results, imc.buildDiscogsRelease(release))
				}
			}
		}
	}
	if len(results) == 0 {
		return nil, errors.New("no releases found for the specified format")
	}
	imc.printState()
	return results, nil
}

func (imc *InMemoryClient) FetchReleasesByYear(yearStr string) ([]models.DiscogsRelease, error) {
	imc.mu.Lock()
	defer imc.mu.Unlock()

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return nil, errors.New("invalid year format")
	}

	var results []models.DiscogsRelease
	for _, release := range imc.releases {
		if release.Year == year {
			results = append(results, imc.buildDiscogsRelease(release))
		}
	}
	if len(results) == 0 {
		return nil, errors.New("no releases found for the specified year")
	}
	imc.printState()
	return results, nil
}

// buildDiscogsRelease constructs a DiscogsRelease from a Release
func (imc *InMemoryClient) buildDiscogsRelease(release models.Release) models.DiscogsRelease {
	var artists []struct {
		Name string `json:"name"`
	}
	for _, artistID := range imc.releasesArtists[release.ID] {
		if artist, exists := imc.artists[artistID]; exists {
			artists = append(artists, struct {
				Name string `json:"name"`
			}{Name: artist.Name})
		}
	}

	var formats []string
	for _, formatID := range imc.releasesFormats[release.ID] {
		if format, exists := imc.formats[formatID]; exists {
			formats = append(formats, format.Name)
		}
	}

	return models.DiscogsRelease{
		ID:          release.ID,
		Title:       release.Title,
		Year:        release.Year,
		Artist:      artists,
		CatalogNo:   release.CatalogNo,
		Thumb:       release.Thumb,
		ResourceURL: release.ResourceURL,
		Formats:     formats,
		Status:      "", // Add status if needed
	}
}

// GetState returns the entire state of the in-memory database
func (imc *InMemoryClient) GetState() map[string]interface{} {
	imc.mu.Lock()
	defer imc.mu.Unlock()

	state := map[string]interface{}{
		"releases": imc.releases,
		"artists":  imc.artists,
		"formats":  imc.formats,
	}
	return state
}

// printState prints the current state of the in-memory database
func (imc *InMemoryClient) printState() {
	fmt.Println("Current state of the in-memory database:")
	fmt.Println("Releases:")
	for _, release := range imc.releases {
		fmt.Printf("ID: %d, Title: %s, Year: %d\n", release.ID, release.Title, release.Year)
	}
	fmt.Println("Artists:")
	for _, artist := range imc.artists {
		fmt.Printf("ID: %d, Name: %s\n", artist.ID, artist.Name)
	}
	fmt.Println("Formats:")
	for _, format := range imc.formats {
		fmt.Printf("ID: %d, Name: %s\n", format.ID, format.Name)
	}
}
