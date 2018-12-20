package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path"
)

// APIInfoHandler returns some information about the API
func APIInfoHandler(a *AppContext) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		type APIInfo struct {
			Version string `json:"version"`
		}

		info := APIInfo{
			Version: version,
		}

		j, err := json.Marshal(info)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(FormatError("The server encountered an error trying to " +
				"encode the JSON response."))
			return
		}

		w.Write(j)

		return

	})

	return handler
}

// FileServer serves static files. It can be set to serve a particular file in place of a
// 404 message in the configuration file. By default, it will serve the 404.
func FileServer(a *AppContext) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if a.Config.DefaultFile != "" {
			p := r.URL.Path
			if _, err := os.Stat(path.Join(a.Config.Webroot, p)); !os.IsNotExist(err) {
				http.FileServer(http.Dir(a.Config.Webroot)).ServeHTTP(w, r)
			} else {
				http.ServeFile(w, r, path.Join(a.Config.Webroot, a.Config.DefaultFile))
			}
			return
		}
		http.FileServer(http.Dir(a.Config.Webroot)).ServeHTTP(w, r)

	})
}
