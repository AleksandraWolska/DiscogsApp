package apiclient

import (
	"net/http"
)

type DiscogsClient struct {
	httpClient *http.Client
	baseURL    string
	token      string // if using OAuth or token-based auth
}

func NewDiscogsClient(httpClient *http.Client, baseURL, token string) *DiscogsClient {
	return &DiscogsClient{
		httpClient: httpClient,
		baseURL:    baseURL,
		token:      token,
	}
}
