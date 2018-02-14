package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/blevesearch/bleve"
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	bf "gopkg.in/russross/blackfriday.v2"
)

// GetPageHandler returns the contents of a page - Markdown or HTML
func GetPageHandler(a *AppContext) (handler http.HandlerFunc) {
	handler = func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		path := a.Config.WikiDir + vars["slug"] + ".md"

		if _, err := os.Stat(path); os.IsNotExist(err) {
			w.WriteHeader(404)
			return
		}

		page, err := ParsePageFile(path)
		if err != nil {
			w.WriteHeader(500)
			log.Println(err)
			return
		}

		format, ok := r.URL.Query()["format"]
		if !ok || len(format) < 1 {
			format = []string{"source"}
		}

		switch format[0] {
		case "html":
			unsafe := bf.Run(bf.Run([]byte(page.Contents)))
			page.Contents = string(bluemonday.UGCPolicy().SanitizeBytes(unsafe))
		default:
			// Do nothing
		}

		j, err := json.Marshal(page)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.Write(j)

	}
	return
}

// PostPageHandler updates the contents of a page - creating it if it doesn't
// exist.
func PostPageHandler(a *AppContext) (handler http.HandlerFunc) {
	handler = func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		path := a.Config.WikiDir + vars["slug"] + ".md"

		decoder := json.NewDecoder(r.Body)
		var p Page
		err := decoder.Decode(&p)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		headerBuf := new(bytes.Buffer)
		if err = toml.NewEncoder(headerBuf).Encode(p.Header()); err != nil {
			w.WriteHeader(400)
			return
		}

		fileContents := fmt.Sprintf("<!-- Devpad Header\n%s-->\n\n%s", headerBuf.String(), p.Contents)

		err = ioutil.WriteFile(path, []byte(fileContents), 0644)
		if err != nil {
			fmt.Printf("Error Saving %s.\n", path)
			w.WriteHeader(500)
			return
		}

	}
	return
}

// DeletePageHandler deletes a markdown file.
func DeletePageHandler(a *AppContext) (handler http.HandlerFunc) {
	handler = func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		path := a.Config.WikiDir + vars["slug"] + ".md"

		err := os.Remove(path)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		return

	}
	return
}

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
