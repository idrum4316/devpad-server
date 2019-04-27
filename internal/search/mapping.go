package search

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/mapping"
)

// NewPageMapping creates the Bleve mapping for a page structure
func NewPageMapping() *mapping.IndexMappingImpl {

	// Mapping for english fields
	enFieldMapping := bleve.NewTextFieldMapping()
	enFieldMapping.Analyzer = "en"

	// Mapping for keyword fields
	kwFieldMapping := bleve.NewTextFieldMapping()
	kwFieldMapping.Analyzer = keyword.Name

	// Set mapping for page.metadata
	metadataMapping := bleve.NewDocumentMapping()
	metadataMapping.AddFieldMappingsAt("tags", kwFieldMapping)

	// Set mapping for page
	pageMapping := bleve.NewDocumentMapping()
	pageMapping.AddFieldMappingsAt("contents", enFieldMapping)
	pageMapping.AddSubDocumentMapping("metadata", metadataMapping)

	m := bleve.NewIndexMapping()
	m.DefaultMapping = pageMapping

	return m

}
