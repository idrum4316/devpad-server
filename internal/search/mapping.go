package search

import (
	"github.com/blevesearch/bleve"
	_ "github.com/blevesearch/bleve/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/mapping"
)

// NewPageMapping creates the Bleve mapping for a page structure
func NewPageMapping() *mapping.IndexMappingImpl {

	// Mapping for english fields
	enFieldMapping := bleve.NewTextFieldMapping()
	enFieldMapping.Analyzer = "en"

	// Mapping for keyword fields
	kwFieldMapping := bleve.NewTextFieldMapping()
	kwFieldMapping.Analyzer = "keyword"

	pageMapping := bleve.NewDocumentMapping()
	pageMapping.AddFieldMappingsAt("contents", enFieldMapping)
	pageMapping.AddFieldMappingsAt("metadata.tags", kwFieldMapping)

	m := bleve.NewIndexMapping()
	m.DefaultMapping = pageMapping

	return m

}