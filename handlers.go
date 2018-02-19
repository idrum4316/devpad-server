package main

import (
	"encoding/json"
	"net/http"
)

// APIInfoHandler returns some information about the API
func APIInfoHandler(a *AppContext) (handler http.HandlerFunc) {
	handler = func(w http.ResponseWriter, r *http.Request) {

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

	}
	return
}
