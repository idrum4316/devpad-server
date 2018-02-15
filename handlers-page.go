package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	bf "gopkg.in/russross/blackfriday.v2"
)

// GetPageHandler returns the contents of a page - Markdown or HTML
func GetPageHandler(a *AppContext) (handler http.HandlerFunc) {
	handler = func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		path := fmt.Sprintf("%s.md", path.Join(a.Config.WikiDir, vars["slug"]))

		if _, err := os.Stat(path); os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
			w.Write(FormatError("The page you requested could not be found."))
			return
		}

		page, err := NewPageFromFile(path)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(FormatError("The server encountered an error trying to " +
				"parse the requested file."))
			return
		}

		// format should be "html" or "source"
		format, ok := r.URL.Query()["format"]
		if !ok || len(format) < 1 {
			format = []string{"source"}
		}

		switch format[0] {
		case "html":
			unsafe := bf.Run([]byte(page.Contents))
			page.Contents = string(bluemonday.UGCPolicy().SanitizeBytes(unsafe))
		case "source":
			// Don't render the Markdown
		default:
			w.WriteHeader(http.StatusBadRequest)
			w.Write(FormatError("Unknown value in 'format' parameter."))
			return
		}

		j, err := json.Marshal(page)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(FormatError("An error occurred occurred trying to format " +
				"a response."))
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
		path := fmt.Sprintf("%s.md", path.Join(a.Config.WikiDir, vars["slug"]))

		decoder := json.NewDecoder(r.Body)
		var p Page
		err := decoder.Decode(&p)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(FormatError("Unable to decode JSON request."))
			return
		}

		headerBuf := new(bytes.Buffer)
		if err = toml.NewEncoder(headerBuf).Encode(p.Header()); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(FormatError("Unable to encode page header."))
			return
		}

		fileContents := fmt.Sprintf("<!-- Devpad Header\n%s-->\n\n%s",
			headerBuf.String(), p.Contents)

		err = ioutil.WriteFile(path, []byte(fileContents), 0644)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(FormatError("Unable to write file to disk."))
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
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(FormatError("Unable to delete file from disk."))
			return
		}

		return

	}
	return
}
