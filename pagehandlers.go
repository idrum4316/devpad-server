package main

import (
	"net/http"
)

// WebpageHandler displays a wiki page's contents
func WebpageHandler(a *AppContext) (handler http.HandlerFunc) {
	handler = func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/app/index.html")
	}
	return
}
