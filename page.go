package main

import (
	"bufio"
	"github.com/BurntSushi/toml"
	"os"
)

type Page struct {
	Title    string   `json:"title"`
	Tags     []string `json:"tags"`
	Contents string   `json:"contents"`
}

type PageHeader struct {
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
}

// NewPage generates a new empty page instance
func NewPage() *Page {
	return &Page{
		Tags: []string{},
	}
}

// Get just the page header
func (p *Page) Header() *PageHeader {
	return &PageHeader{
		Title: p.Title,
		Tags:  p.Tags,
	}
}

// Expand a PageHeader to a Page with an empty 'Contents'
func (h *PageHeader) ToPage() *Page {
	return &Page{
		Title: h.Title,
		Tags:  h.Tags,
	}
}

// ParseHeader parses the page header from the markdown, returning the markdown
// without the header.
func ParsePageFile(path string) (page *Page, err error) {

	page = NewPage()

	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

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
