package models

type Pagination struct {
	Pages int `json:"pages"`
}
type DiscogsApiResponse struct {
	Releases   []DiscogsApiRelease `json:"releases"`
	Pagination Pagination          `json:"pagination"`
}

// Representation of a discogs release
type DiscogsApiRelease struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Year         int    `json:"year"`
	Artist       string `json:"artist"`
	Status       string `json:"status"`
	Format       string `json:"format"`
	CatalogNo    string `json:"catno"`
	ThumbnailURL string `json:"thumb"`
	ResourceURL  string `json:"resource_url"`
}
