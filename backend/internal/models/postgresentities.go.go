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
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Format struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Status struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// etc.

// Representation of a discogs release
type DiscogsRelease struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Year        int      `json:"year"`
	Artists     []Artist `json:"artists"`
	CatalogNo   string   `json:"catno"`
	Thumb       string   `json:"thumb"`
	ResourceURL string   `json:"resource_url"`
	Formats     []Format `json:"formats"`
	Status      string   `json:"status"`
}
