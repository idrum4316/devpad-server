package main

import (
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// A search result is an array of files along with the number of occurrences
// found in that file
type searchResult struct {
	Slug        string `json:"slug"`
	Occurrences int    `json:"count"`
}

// Search all files in the wiki directory for the search term
func searchWiki(query string, directory string) (result []searchResult, err error) {
	result = []searchResult{}

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Println(err)
		return
	}

	re := regexp.MustCompile("(?i)" + query)

	for _, f := range files {
		if !f.IsDir() {
			contents, err := ioutil.ReadFile(path.Join(directory, f.Name()))
			if err != nil {
				log.Println(err)
				return result, err
			}

			num := len(re.FindAllString(string(contents), -1))
			if num > 0 {
				result = append(result, searchResult{strings.TrimSuffix(f.Name(), filepath.Ext(f.Name())), num})
			}
		}
	}

	return result, nil
}
