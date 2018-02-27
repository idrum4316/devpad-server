package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/blevesearch/bleve"
	"github.com/fsnotify/fsnotify"
)

// A search result is an array of files along with the number of occurrences
// found in that file
type searchResult struct {
	Slug        string `json:"slug"`
	Occurrences int    `json:"count"`
}

// Initialize the search index
func initSearchIndex(a *AppContext) error {
	mapping := NewPageMapping()

	if a.Config.IndexInMemory {
		index, err := bleve.NewMemOnly(mapping)
		if err != nil {
			log.Fatal(err)
		}
		a.SearchIndex = index
	} else {

		if a.Config.IndexLoc != "" {
			os.RemoveAll(a.Config.IndexLoc)
		}

		index, err := bleve.New(a.Config.IndexLoc, mapping)
		if err != nil {
			log.Fatal(err)
		}
		a.SearchIndex = index

	}

	return nil
}

// Index all markdown files in the wiki folder
func indexAll(c *AppContext) {
	files, err := ioutil.ReadDir(c.Config.WikiDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() {
			continue // Skip directories
		}
		if !strings.HasSuffix(f.Name(), ".md") {
			continue // Skip non markdown files
		}

		page, err := NewPageFromFile(path.Join(c.Config.WikiDir, f.Name()))
		if err != nil {
			log.Fatal(err)
		}

		c.SearchIndex.Index(strings.TrimSuffix(f.Name(), filepath.Ext(f.Name())), page)
	}

	watchForChanges(c)
	return
}

// Index all files that change
func watchForChanges(c *AppContext) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				// Only index markdown files
				if !strings.HasSuffix(event.Name, ".md") {
					continue
				}
				updateSearchIndex(c, event)
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(c.Config.WikiDir)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

// Update a single file in the search index
func updateSearchIndex(c *AppContext, e fsnotify.Event) {
	_, file := path.Split(e.Name)

	// Handle each file event
	switch e.Op {

	case fsnotify.Write:
		page, err := NewPageFromFile(path.Join(c.Config.WikiDir, file))
		if err != nil {
			log.Fatal(err)
		}
		c.SearchIndex.Index(strings.TrimSuffix(file, filepath.Ext(file)), page)

	case fsnotify.Create:
		page, err := NewPageFromFile(path.Join(c.Config.WikiDir, file))
		if err != nil {
			log.Fatal(err)
		}
		c.SearchIndex.Index(strings.TrimSuffix(file, filepath.Ext(file)), page)

	case fsnotify.Remove:
		err := c.SearchIndex.Delete(strings.TrimSuffix(file, filepath.Ext(file)))
		if err != nil {
			log.Fatal(err)
		}

	case fsnotify.Rename:
		if _, err := os.Stat(path.Join(c.Config.WikiDir, file)); os.IsNotExist(err) {
			err := c.SearchIndex.Delete(strings.TrimSuffix(file, filepath.Ext(file)))
			if err != nil {
				log.Fatal(err)
			}
		}

	}

}
