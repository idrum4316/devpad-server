package search

import (
	"github.com/microcosm-cc/bluemonday"
)

var htmlPolicy *bluemonday.Policy

func init() {
	htmlPolicy = bluemonday.StrictPolicy()
}
