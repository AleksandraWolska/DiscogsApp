package models

type Release struct {
	ID          int
	Title       string
	Year        int
	CatalogNo   string
	Thumb       string
	ResourceURL string
	// ...
}

type Artist struct {
	ID   int
	Name string `json:"name"`
}

type Format struct {
	ID   int
	Name string
}

type Status struct {
	ID   int
	Name string `json:"name"`
}

// etc.

// Representation of a discogs release
type DiscogsRelease struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Year        int    `json:"year"`
	Artists     []Artist
	CatalogNo   string `json:"catno"`
	Thumb       string `json:"thumb"`
	ResourceURL string `json:"resource_url"`
	Formats     []Format
	Status      string `json:"status"`
}
