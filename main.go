package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {

	// Load the configuration file
	appContext := NewAppContext()
	appContext.Config.LoadFromFile("config.toml")
	err := initSearchIndex(appContext)
	if err != nil {
		log.Fatal(err)
	}
	go indexAll(appContext)

	// Create the wiki directory if it doesn't exist
	if _, err = os.Stat(appContext.Config.WikiDir); os.IsNotExist(err) {
		os.MkdirAll(appContext.Config.WikiDir, 0775)
	}

	// HTTP Router
	router := mux.NewRouter()

	// Handle the API calls
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/pages/{slug}", GetPageHandler(appContext)).Methods("GET")
	apiRouter.HandleFunc("/pages/{slug}", PostPageHandler(appContext)).Methods("POST")
	apiRouter.HandleFunc("/pages/{slug}", DeletePageHandler(appContext)).Methods("DELETE")
	apiRouter.HandleFunc("/search", SearchHandler(appContext)).Methods("GET")

	// Serves static files
	if appContext.Config.ServeStatic {
		webroot := appContext.Config.Webroot
		router.PathPrefix("/").Handler(http.FileServer(http.Dir(webroot)))
	}

	// Start the server
	logged_router := handlers.LoggingHandler(os.Stdout, router)
	host := appContext.Config.ListenHost
	port := appContext.Config.Port
	err = http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), logged_router)
	if err != nil {
		fmt.Println(err)
	}

}
