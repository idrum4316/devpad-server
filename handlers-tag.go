package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/blevesearch/bleve"
)

// GetTagsHandler returns a list of all tags
func GetTagsHandler(a *AppContext) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		numTags := 10000

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
			numTags = sizeInt
		}

		query := bleve.NewMatchAllQuery()
		search := bleve.NewSearchRequest(query)
		search.Size = 0
		tagsFacet := bleve.NewFacetRequest("metadata.tags", numTags)
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

	})

	return RequireAuth(handler, a)
}
