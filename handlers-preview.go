package main

import (
	"encoding/json"
	"net/http"

	"github.com/Depado/bfchroma"
	"github.com/microcosm-cc/bluemonday"
	bf "gopkg.in/russross/blackfriday.v2"
)

// Render the text sent in to HTML
func PostPreviewHandler(a *AppContext) (handler http.HandlerFunc) {
	handler = func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		page := NewPage()
		err := decoder.Decode(&page)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(FormatError("Unable to decode JSON request."))
			return
		}

		// toc should be "true" or "false"
		toc, ok := r.URL.Query()["toc"]
		if !ok || len(toc) < 1 {
			toc = []string{"false"}
		}

		renderer := bf.NewHTMLRenderer(bf.HTMLRendererParameters{
			Flags: bf.CommonHTMLFlags,
		})

		if toc[0] == "true" {
			renderer.Flags |= bf.TOC
		}

		if a.Config.SanitizeHTML {
			unsafe := bf.Run([]byte(page.Contents), bf.WithRenderer(renderer))
			page.Contents = string(bluemonday.UGCPolicy().SanitizeBytes(unsafe))
		} else {
			r := bfchroma.NewRenderer(
				bfchroma.Extend(renderer),
				bfchroma.WithoutAutodetect(),
				bfchroma.Style("tango"),
			)
			page.Contents = string(bf.Run([]byte(page.Contents), bf.WithRenderer(r)))
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
