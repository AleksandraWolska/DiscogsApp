package apiclient

import (
	"discogsbackend/internal/models"
	"net/http"
)

type DiscogsClient struct {
	httpClient *http.Client
	baseURL    string
}

type DiscogsClientInterface interface {
	FetchOnLabel(labelID string) ([]models.DiscogsApiRelease, error)
}

func NewDiscogsClient(httpClient *http.Client, baseURL, token string) *DiscogsClient {
	return &DiscogsClient{
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}
