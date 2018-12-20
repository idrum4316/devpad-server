package main

import (
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

// RequireAuth checks for a valid token before forwarding
func RequireAuth(next http.Handler, a *AppContext) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenHeader := r.Header.Get("jwt")
		if tokenHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(FormatError("unauthorized"))
			return
		}

		token, err := jwt.Parse(tokenHeader, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(a.Config.SigningKey), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(FormatError("error reading token"))
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userExists, err := a.Store.UserExists(claims["userid"].(string))
			if !userExists || err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write(FormatError("unauthorized"))
				return
			}

			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(FormatError("unauthorized"))
			return
		}

	})
}
