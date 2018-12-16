package search

import (
	"os"

	"github.com/blevesearch/bleve"
	"github.com/idrum4316/devpad-server/internal/page"
)

// Index is the search index
type Index struct {
	index bleve.Index
}

// NewIndex returns a new Index instance
func NewIndex(path string) (*Index, error) {

	var index bleve.Index
	var err error

	if _, err := os.Stat(path); os.IsNotExist(err) {
		mapping := NewPageMapping()
		index, err = bleve.New(path, mapping)
	} else {
		index, err = bleve.Open(path)
	}

	if err != nil {
		return nil, err
	}

	i := Index{
		index: index,
	}

	return &i, nil

}

// Close the bleve database
func (i *Index) Close() error {
	err := i.index.Close()
	return err
}

// ExecuteSearch executes the search query in the index and returns the result
func (i *Index) ExecuteSearch(r *bleve.SearchRequest) (*bleve.SearchResult, error) {

	result, err := i.index.Search(r)
	return result, err

}

// IndexPage adds or updates the page in the index
func (i *Index) IndexPage(id string, p *page.Page) error {

	// Remove all html tags from the page before indexing
	p.Contents = htmlPolicy.Sanitize(p.Contents)

	err := i.index.Index(id, p)
	return err

}
