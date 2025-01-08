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
	Name string
}

type Format struct {
	ID   int
	Name string
}

type Status struct {
	ID   int
	Name string
}

// etc.

// Representation of a discogs release
type DiscogsRelease struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Year   int    `json:"year"`
	Artist []struct {
		Name string `json:"name"`
	} `json:"artists"`
	CatalogNo   string   `json:"catno"`
	Thumb       string   `json:"thumb"`
	ResourceURL string   `json:"resource_url"`
	Formats     []string `json:"formats"`
	Status      string   `json:"status"`
}
