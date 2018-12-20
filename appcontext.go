package main

import (
	"errors"
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
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

// GetUserIDFromRequest returns the user id from the JWT token in the request
func (a *AppContext) GetUserIDFromRequest(r *http.Request) (id string, err error) {

	// Get the token from the HTTP headers. If it doesn't exist, return an
	// error.
	tokenHeader := r.Header.Get("jwt")
	if tokenHeader == "" {
		err = errors.New("missing jwt token")
		return
	}

	// Parse the token. Return any errors that are thrown.
	token, err := jwt.Parse(tokenHeader, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.Config.SigningKey), nil
	})
	if err != nil {
		return
	}

	// If the token is good, return the user id. If it's not, return an error.
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id = claims["userid"].(string)
		return
	}

	// If this code runs, that means a valid claim was not found
	err = errors.New("invalid token")
	return

}
