package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/idrum4316/devpad-server/internal/datastore"
	"github.com/idrum4316/devpad-server/internal/search"
)

var version = "0.0.6"

func main() {

	// Load the configuration file
	appContext := NewAppContext()
	appContext.Config.LoadFromFile("config.toml")

	// Create the data directory if it doesn't exist
	if _, err := os.Stat(appContext.Config.DataDir); os.IsNotExist(err) {
		os.MkdirAll(appContext.Config.DataDir, 0700)
	}

	// Create and attach the Bolt datastore
	store, err := datastore.New(path.Join(appContext.Config.DataDir, "devpad.db"))
	if err != nil {
		log.Fatal(err)
	}
	appContext.Store = store
	defer appContext.Store.Close()

	// Create and attach the Bleve search index
	index, err := search.NewIndex(path.Join(appContext.Config.DataDir, "pages.index"))
	if err != nil {
		log.Fatal(err)
	}
	appContext.Index = index
	defer appContext.Index.Close()

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
	apiRouter.HandleFunc("/preview", PostPreviewHandler(appContext)).Methods("POST")

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
