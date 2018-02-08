package main

import (
	"github.com/BurntSushi/toml"
)

type AppConfig struct {
	Port        int
	ListenHost  string
	WikiDir     string
	ServeStatic bool
	Webroot     string
}

func NewAppConfig() (c *AppConfig) {
	c = &AppConfig{
		Port:        8080,
		ListenHost:  "127.0.0.1",
		WikiDir:     "./wiki/",
		ServeStatic: true,
		Webroot:     "./wwwroot/",
	}
	return
}

// LoadFromFile loads the toml config from the configuration file.
func (c *AppConfig) LoadFromFile(file string) (err error) {
	_, err = toml.DecodeFile(file, c)
	return
}
