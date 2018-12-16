package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
)

// SearchHandler searches the wiki files for a search term
func SearchHandler(a *AppContext) (handler http.HandlerFunc) {
	handler = func(w http.ResponseWriter, r *http.Request) {

		searchQuery := ""

		searchInput, ok := r.URL.Query()["q"]
		if ok {
			if len(searchInput) > 0 {
				searchQuery = searchInput[0]
			}
		}

		queries := []query.Query{}

		if searchQuery == "" {
			queries = append(queries, bleve.NewMatchAllQuery())
		} else {
			queries = append(queries, bleve.NewQueryStringQuery(searchQuery))
		}

		// Check for the 'size' parameter
		tags, ok := r.URL.Query()["tag"]
		if ok {
			for _, tag := range tags {
				tagQuery := bleve.NewTermQuery(tag)
				tagQuery.FieldVal = "metadata.tags"
				queries = append(queries, tagQuery)
			}
		}

		q := bleve.NewConjunctionQuery(queries...)
		search := bleve.NewSearchRequest(q)
		search.Highlight = bleve.NewHighlight()
		search.Fields = []string{"contents", "metadata.title", "metadata.tags", "metadata.modified"}

		// Check for the 'size' parameter
		size, ok := r.URL.Query()["size"]
		if ok {
			sizeInt, err := strconv.Atoi(size[0])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(FormatError("Unable to parse integer from 'size'" +
					" option."))
				return
			}
			search.Size = sizeInt
		}

		// Check for the 'from' parameter
		from, ok := r.URL.Query()["from"]
		if ok {
			fromInt, err := strconv.Atoi(from[0])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(FormatError("Unable to parse integer from 'from'" +
					" option."))
				return
			}
			search.From = fromInt
		}

		// Check for the 'sort' paramter
		sort, ok := r.URL.Query()["sort"]
		fmt.Printf("Sort By: %+v\n", sort)
		if ok {
			search.SortBy(sort)
		}

		// Add the Tags facet
		tagsFacet := bleve.NewFacetRequest("metadata.tags", 100)
		search.AddFacet("tags", tagsFacet)

		searchResults, err := a.Index.ExecuteSearch(search)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(FormatError("Unable to process your search query."))
			return
		}

		j, err := json.Marshal(searchResults)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(FormatError("Unable to encode the response."))
			return
		}
		w.Write(j)

	}
	return
}
