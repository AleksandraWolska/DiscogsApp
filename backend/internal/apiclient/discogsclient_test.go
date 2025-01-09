package apiclient

import (
	"discogsbackend/internal/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFetchOnLabel_Success(t *testing.T) {
	labelID := "123"
	expectedResponse := models.DiscogsApiResponse{
		Pagination: models.Pagination{Pages: 1},
		Releases: []models.DiscogsApiRelease{
			{ID: 1, Title: "Release 1"},
			{ID: 2, Title: "Release 2"},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	}))
	defer server.Close()

	client := &DiscogsClient{
		baseURL:    server.URL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}

	releases, err := client.FetchOnLabel(labelID)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(releases) != len(expectedResponse.Releases) {
		t.Errorf("expected %d releases but got %d", len(expectedResponse.Releases), len(releases))
	}

	for i, release := range releases {
		if release.ID != expectedResponse.Releases[i].ID || release.Title != expectedResponse.Releases[i].Title {
			t.Errorf("expected release %v but got %v", expectedResponse.Releases[i], release)
		}
	}
}
