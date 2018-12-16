package main

import (
	"github.com/idrum4316/devpad-server/internal/datastore"
	"github.com/idrum4316/devpad-server/internal/search"
)

// AppContext holds the overall application context (config, etc..)
type AppContext struct {
	Config *AppConfig
	Index  *search.Index
	Store  *datastore.Datastore
}

// NewAppContext returns a pointer to a new AppContext with default values set.
func NewAppContext() (a *AppContext) {
	a = &AppContext{
		Config: NewAppConfig(),
		Index:  nil,
		Store:  nil,
	}
	return
}
