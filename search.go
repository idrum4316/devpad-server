package main

import (
	"github.com/blevesearch/bleve"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// A search result is an array of files along with the number of occurrences
// found in that file
type searchResult struct {
	Slug        string `json:"slug"`
	Occurrences int    `json:"count"`
}

// Search all files in the wiki directory for the search term
func searchWiki(c *AppContext, q string) (result *bleve.SearchResult, err error) {
	query := bleve.NewMatchQuery(q)
	search := bleve.NewSearchRequest(query)
	search.Highlight = bleve.NewHighlight()
	searchResults, err := c.SearchIndex.Search(search)
	if err != nil {
		return nil, err
	}
	return searchResults, nil
}

// Initialize the search index
func initSearchIndex(a *AppContext) error {
	mapping := bleve.NewIndexMapping()

	if a.Config.IndexInMemory {
		index, err := bleve.NewMemOnly(mapping)
		if err != nil {
			log.Fatal(err)
		}
		a.SearchIndex = index
	} else {
		if _, err := os.Stat(a.Config.IndexFile); err == nil {
			index, err := bleve.Open(a.Config.IndexFile)
			if err != nil {
				log.Fatal(err)
			}
			a.SearchIndex = index
		} else {
			index, err := bleve.New(a.Config.IndexFile, mapping)
			if err != nil {
				log.Fatal(err)
			}
			a.SearchIndex = index
		}
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
			// Skip directories
			continue
		}
		if !strings.HasSuffix(f.Name(), ".md") {
			// Skip non markdown files
			continue
		}

		page, err := ParsePageFile(path.Join(c.Config.WikiDir, f.Name()))
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

	switch e.Op {
	case fsnotify.Write:
		page, err := ParsePageFile(path.Join(c.Config.WikiDir, file))
		if err != nil {
			log.Fatal(err)
		}
		c.SearchIndex.Index(strings.TrimSuffix(file, filepath.Ext(file)), page)

	case fsnotify.Create:
		page, err := ParsePageFile(path.Join(c.Config.WikiDir, file))
		if err != nil {
			log.Fatal(err)
		}
		c.SearchIndex.Index(strings.TrimSuffix(file, filepath.Ext(file)), page)

	case fsnotify.Remove:
		err := c.SearchIndex.Delete(strings.TrimSuffix(file, filepath.Ext(file)))
		if err != nil {
			log.Fatal(err)
		}

	}
}
