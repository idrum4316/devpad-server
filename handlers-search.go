package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/blevesearch/bleve"
)

// SearchHandler searches the wiki files for a search term
func SearchHandler(a *AppContext) (handler http.HandlerFunc) {
	handler = func(w http.ResponseWriter, r *http.Request) {

		queryString := r.URL.Query()["q"]
		searchQuery := ""

		if len(queryString) > 0 {
			searchQuery = queryString[0]
		}

		query := bleve.NewQueryStringQuery(searchQuery)
		search := bleve.NewSearchRequest(query)
		search.Highlight = bleve.NewHighlight()
		search.Fields = []string{"title"}

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

		searchResults, err := a.SearchIndex.Search(search)
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
