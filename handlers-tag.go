package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/blevesearch/bleve"
	"github.com/gorilla/mux"
)

// GetTagsHandler returns a list of all tags
func GetTagsHandler(a *AppContext) (handler http.HandlerFunc) {
	handler = func(w http.ResponseWriter, r *http.Request) {

		query := bleve.NewMatchAllQuery()
		search := bleve.NewSearchRequest(query)
		search.Size = 0
		tagsFacet := bleve.NewFacetRequest("tags", 10000)
		search.AddFacet("tags", tagsFacet)
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

// Return the pages with a specific tag
func GetTagHandler(a *AppContext) (handler http.HandlerFunc) {
	handler = func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		query := bleve.NewQueryStringQuery(fmt.Sprintf("tags:\"%s\"", vars["tag"]))
		search := bleve.NewSearchRequest(query)
		search.Fields = []string{"title"}
		search.Size = 10000
		search.SortBy([]string{"title"})
		searchResults, err := a.SearchIndex.Search(search)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(FormatError("An error occurred while running the search."))
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
