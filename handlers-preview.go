package main

import (
	"encoding/json"
	"net/http"

	"github.com/Depado/bfchroma"
	"github.com/idrum4316/devpad-server/internal/page"
	"github.com/microcosm-cc/bluemonday"
	bf "gopkg.in/russross/blackfriday.v2"
)

// PostPreviewHandler renders the text sent in to HTML
func PostPreviewHandler(a *AppContext) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		pg := page.Page{}
		err := decoder.Decode(&pg)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(FormatError("Unable to decode JSON request."))
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
