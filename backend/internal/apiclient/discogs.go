package apiclient

import (
	"discogsbackend/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
)

// FetchOnLabel fetches all releases for a given label ID or label name
func (dc *DiscogsClient) FetchOnLabel(labelID string) ([]models.DiscogsApiRelease, error) {
	var allReleases []models.DiscogsApiRelease
	page := 1

	for {
		url := fmt.Sprintf("%s/labels/%s/releases?page=%d", dc.baseURL, labelID, page)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", "MyDiscogsApp/1.0 +http://myapp.example.com")

		resp, err := dc.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var result models.DiscogsApiResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, err
		}

		allReleases = append(allReleases, result.Releases...)

		if page >= result.Pagination.Pages {
			break
		}
		page++
	}

	return allReleases, nil
}
