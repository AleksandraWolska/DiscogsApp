package inmemorydatabase

import (
	"discogsbackend/internal/models"
	"sync"
)

// InMemoryClient holds the in-memory data
type InMemoryClient struct {
	releases        map[int]models.Release
	artists         map[int]models.Artist
	formats         map[int]models.Format
	statuses        map[int]models.Status
	releasesArtists map[int][]int // Map to store which artists are associated with which releases
	releasesFormats map[int][]int // Map to store which formats are associated with which releases
	nextID          int
	mu              sync.Mutex
}

// NewInMemoryClient creates a new InMemoryClient
func NewInMemoryClient() *InMemoryClient {
	return &InMemoryClient{
		releases:        make(map[int]models.Release),
		artists:         make(map[int]models.Artist),
		formats:         make(map[int]models.Format),
		statuses:        make(map[int]models.Status),
		releasesArtists: make(map[int][]int),
		releasesFormats: make(map[int][]int),
		nextID:          1,
	}
}
