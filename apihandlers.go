package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	bf "gopkg.in/russross/blackfriday.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
		var d PageData
		err := decoder.Decode(&d)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		fileContents := fmt.Sprintf(`<!-- TinyWiki Header
title = "%s"
-->

%s`, d.Title, d.Contents)

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

		results, err := searchWiki(searchQuery, a.Config.WikiDir)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		j, err := json.Marshal(results)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.Write(j)

	}
	return
}
