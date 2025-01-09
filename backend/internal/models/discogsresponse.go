package models

type DiscogsApiResponse struct {
	Releases   []DiscogsApiRelease `json:"releases"`
	Pagination struct {
		Pages int `json:"pages"`
	} `json:"pagination"`
}

// type Release struct {
// 	ID    int
// 	Title string
// 	Year  int
// 	// ...
// }

// type Artist struct {
// 	ID   int
// 	Name string
// }

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
