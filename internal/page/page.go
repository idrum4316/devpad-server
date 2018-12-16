package page

import (
	"time"
)

// Page represents a single wiki page
type Page struct {
	Contents string   `json:"contents"`
	Metadata Metadata `json:"metadata"`
}

// Metadata contains the metadata for the page
type Metadata struct {
	Title    string    `json:"title"`
	Tags     []string  `json:"tags"`
	Modified time.Time `json:"modified"`
}

// New generates a new empty page instance
func New() *Page {
	return &Page{
		Metadata: Metadata{
			Tags: []string{},
		},
	}
}
