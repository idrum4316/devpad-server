package main

// PageHeader contains a subset of the fields of a Page. These are written as
// the header at the top of the markdown files.
type PageHeader struct {
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
}

// Expand a PageHeader to a Page
func (h *PageHeader) ToPage() *Page {
	return &Page{
		Title: h.Title,
		Tags:  h.Tags,
	}
}
