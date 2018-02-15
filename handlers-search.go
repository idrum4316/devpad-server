package main

import (
	"encoding/json"
	"net/http"

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
		search.Size = 10000
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
