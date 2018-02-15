package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/blevesearch/bleve"
	"github.com/gorilla/mux"
)

// Return the pages with a specific tag
func TagHandler(a *AppContext) (handler http.HandlerFunc) {
	handler = func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		query := bleve.NewQueryStringQuery(fmt.Sprintf("tags:\"%s\"", vars["tag"]))
		search := bleve.NewSearchRequest(query)
		search.Fields = []string{"title"}
		search.Size = 10000
		search.SortBy([]string{"title"})
		searchResults, err := a.SearchIndex.Search(search)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		j, err := json.Marshal(searchResults)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.Write(j)

	}
	return
}
