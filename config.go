package main

import (
	"github.com/BurntSushi/toml"
)

// AppConfig is the main configuration struct of the devpad-server application.
type AppConfig struct {
	Port          int
	ListenHost    string
	WikiDir       string
	ServeStatic   bool
	Webroot       string
	IndexInMemory bool
	IndexFile     string
}

// NewAppConfig is a constructor that returns a new AppConfig instance with some
// default values set.
func NewAppConfig() (c *AppConfig) {
	c = &AppConfig{
		Port:          8080,
		ListenHost:    "127.0.0.1",
		WikiDir:       "./wiki/",
		ServeStatic:   true,
		Webroot:       "./wwwroot/",
		IndexInMemory: false,
		IndexFile:     "./pages.index",
	}
	return
}

// LoadFromFile loads the toml config from the configuration file into the
// AppConfig instance.
func (c *AppConfig) LoadFromFile(file string) (err error) {
	_, err = toml.DecodeFile(file, c)
	return
}
