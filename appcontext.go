package main

import (
	"github.com/blevesearch/bleve"
)

// AppContext holds the overall application context (config, etc..)
type AppContext struct {
	Config      *AppConfig
	SearchIndex bleve.Index
}

// Return a pointer to a new AppContext with default values set
func NewAppContext() (a *AppContext) {
	a = &AppContext{
		Config:      NewAppConfig(),
		SearchIndex: nil,
	}
	return
}
