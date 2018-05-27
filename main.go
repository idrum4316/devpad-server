package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var version = "0.0.6"

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
	apiRouter.HandleFunc("", APIInfoHandler(appContext)).Methods("GET")
	apiRouter.HandleFunc("/pages", GetPagesHandler(appContext)).Methods("GET")
	apiRouter.HandleFunc("/pages/{slug}", GetPageHandler(appContext)).Methods("GET")
	apiRouter.HandleFunc("/pages/{slug}", PutPageHandler(appContext)).Methods("PUT")
	apiRouter.HandleFunc("/pages/{slug}", DeletePageHandler(appContext)).Methods("DELETE")
	apiRouter.HandleFunc("/search", SearchHandler(appContext)).Methods("GET")
	apiRouter.HandleFunc("/tags", GetTagsHandler(appContext)).Methods("GET")

	// Serves static files
	if appContext.Config.ServeStatic {
		router.PathPrefix("/").Handler(FileServer(appContext))
	}

	// Start the server
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	host := appContext.Config.ListenHost
	port := appContext.Config.Port
	err = http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), loggedRouter)
	if err != nil {
		fmt.Println(err)
	}

}
