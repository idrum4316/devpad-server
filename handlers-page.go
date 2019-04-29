package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Depado/bfchroma"
	"github.com/blevesearch/bleve"
	"github.com/gorilla/mux"
	"github.com/idrum4316/devpad-server/internal/page"
	"github.com/microcosm-cc/bluemonday"
	bf "gopkg.in/russross/blackfriday.v2"
)

// GetPagesHandler returns a list of all pages - with optional paging and
// sorting
func GetPagesHandler(a *AppContext) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		query := bleve.NewMatchAllQuery()
		search := bleve.NewSearchRequest(query)
		search.Fields = []string{"metadata.title", "metadata.tags", "metadata.modified"}

		// Check for the 'size' parameter
		size, ok := r.URL.Query()["size"]
		if ok {
			sizeInt, err := strconv.Atoi(size[0])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write(FormatError("Unable to parse integer from 'size'" +
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
				_, _ = w.Write(FormatError("Unable to parse integer from 'from'" +
					" option."))
				return
			}
			search.From = fromInt
		}

		// Check for the 'sort' parameter
		sortOrder, ok := r.URL.Query()["sort"]
		if ok {
			search.SortBy(sortOrder)
		}

		searchResults, err := a.Index.ExecuteSearch(search)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(FormatError("Unable to process your search query."))
			return
		}

		j, err := json.Marshal(searchResults)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(FormatError("Unable to encode the response."))
			return
		}
		_, _ = w.Write(j)

	})

	return RequireAuth(handler, a)
}

// GetPageHandler returns the contents of a page - Markdown or HTML
func GetPageHandler(a *AppContext) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		slug := vars["slug"]

		pg, err := a.Store.GetPage(slug)

		// Do this if there was an error loading the page (the page not
		// existing is not an error).
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(FormatError("The server encountered an error trying to " +
				"load the requested page."))
			return
		}

		// Do this if the page doesn't exist
		if pg == nil {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write(FormatError("The page you requested could not be found."))
			return
		}

		// format should be "html" or "source"
		format, ok := r.URL.Query()["format"]
		if !ok || len(format) < 1 {
			format = []string{"source"}
		}

		// toc should be "true" or "false"
		toc, ok := r.URL.Query()["toc"]
		if !ok || len(toc) < 1 {
			toc = []string{"false"}
		}

		switch format[0] {
		case "html":
			renderer := bf.NewHTMLRenderer(bf.HTMLRendererParameters{
				Flags: bf.CommonHTMLFlags,
			})

			if toc[0] == "true" {
				renderer.Flags |= bf.TOC
			}

			if a.Config.SanitizeHTML {
				unsafe := bf.Run([]byte(pg.Contents), bf.WithRenderer(renderer))
				pg.Contents = string(bluemonday.UGCPolicy().SanitizeBytes(unsafe))
			} else {
				r := bfchroma.NewRenderer(
					bfchroma.Extend(renderer),
					bfchroma.WithoutAutodetect(),
					bfchroma.Style("tango"),
				)
				pg.Contents = string(bf.Run([]byte(pg.Contents), bf.WithRenderer(r)))
			}

		case "source":
			// Don't render the Markdown
		default:
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(FormatError("Unknown value in 'format' parameter."))
			return
		}

		j, err := json.Marshal(pg)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(FormatError("An error occurred occurred trying to format " +
				"a response."))
			return
		}
		_, _ = w.Write(j)

	})

	return RequireAuth(handler, a)
}

// PutPageHandler updates the contents of a page - creating it if it doesn't
// exist.
func PutPageHandler(a *AppContext) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		decoder := json.NewDecoder(r.Body)
		pg := page.New()
		err := decoder.Decode(pg)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(FormatError("Unable to decode JSON request."))
			log.Println(err)
			return
		}

		// Update the page in datastore
		err = a.Store.UpdatePage(pg, vars["slug"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(FormatError("Unable to save page."))
			return
		}

		// Update the page in the search index
		err = a.Index.IndexPage(vars["slug"], pg)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(FormatError("Unable to update search index."))
			return
		}

	})

	return RequireAuth(handler, a)
}

// DeletePageHandler deletes a markdown file.
func DeletePageHandler(a *AppContext) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		pageID := vars["slug"]

		err := a.Store.DeletePage(pageID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(FormatError("Unable to delete page."))
			return
		}

		err = a.Index.DeletePage(pageID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(FormatError("Unable to remove page from index."))
			return
		}

		return

	})

	return RequireAuth(handler, a)
}

// RenamePageHandler renames a page in the data store.
func RenamePageHandler(a *AppContext) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		pageID := vars["slug"]
		newID := ""

		// Get the new page id from the query string
		newIDParam, ok := r.URL.Query()["id"]
		if ok {
			if len(newIDParam) != 1 {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write(FormatError("You must provide a single new page ID."))
				return
			} else {
				newID = newIDParam[0]
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(FormatError("You must provide a new page ID."))
			return
		}

		// Rename the page
		err := a.Store.RenamePage(pageID, newID)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(FormatError("Unable to rename page."))
			return
		}

		// Remove the old search index data
		err = a.Index.DeletePage(pageID)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(FormatError("Unable to remove old page from index."))
			return
		}

		// Index the new page ID
		pg, err := a.Store.GetPage(newID)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(FormatError("Unable to update search index (1)."))
			return
		}
		err = a.Index.IndexPage(newID, pg)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(FormatError("Unable to update search index (2)."))
			return
		}

		return

	})

	return RequireAuth(handler, a)
}
