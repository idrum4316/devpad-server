package main

// AppContext holds the overall application context (config, etc..)
type AppContext struct {
	Config *AppConfig
}

// Return a pointer to a new AppContext with default values set
func NewAppContext() (a *AppContext) {
	a = &AppContext{
		Config: NewAppConfig(),
	}
	return
}
