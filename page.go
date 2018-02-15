package main

import (
	"bufio"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

// Page represents a single Markdown file.
type Page struct {
	Title    string    `json:"title"`
	Tags     []string  `json:"tags"`
	Modified time.Time `json:"modified"`
	Contents string    `json:"contents"`
}

// NewPage generates a new empty page instance
func NewPage() *Page {
	return &Page{
		Tags: []string{},
	}
}

// Header converts a Page to a PageHeader. Be aware that this process loses the
// Contents and Modified data.
func (p *Page) Header() *PageHeader {
	return &PageHeader{
		Title: p.Title,
		Tags:  p.Tags,
	}
}

// NewPageFromFile parses the page from the markdown, returning the markdown
// without the header.
func NewPageFromFile(path string) (page *Page, err error) {

	page = NewPage()

	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	// This is necessary to get the modification date
	fileInfo, err := file.Stat()
	if err != nil {
		return
	}
	page.Modified = fileInfo.ModTime()

	header := ""
	contents := ""

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	if scanner.Text() == "<!-- Devpad Header" {
		for scanner.Scan() {
			line := scanner.Text()
			if line == "-->" {
				break
			}
			header = header + line + "\n"
		}

		// There's a blank line between the header and the markdown
		scanner.Scan()
		if scanner.Text() != "" {
			contents = contents + scanner.Text() + "\n"
		}

		for scanner.Scan() {
			contents = contents + scanner.Text() + "\n"
		}
	} else {
		contents = scanner.Text()
		for scanner.Scan() {
			contents = contents + scanner.Text() + "\n"
		}
	}

	if _, err = toml.Decode(header, &page); err != nil {
		return
	}

	page.Contents = contents

	return

}
