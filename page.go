package main

import (
	"bufio"
	"github.com/BurntSushi/toml"
	"os"
)

type PageHeader struct {
	Title string
	Tags  []string
}

type Page struct {
	Header   PageHeader
	Contents string
}

// ParseHeader parses the page header from the markdown, returning the markdown
// without the header.
func ParsePageFile(path string) (page *Page, err error) {

	page = &Page{
		Header: PageHeader{
			Tags: []string{},
		},
	}

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

	if _, err = toml.Decode(header, &page.Header); err != nil {
		return
	}

	page.Contents = contents

	return

}
