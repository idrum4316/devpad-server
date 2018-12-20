package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

// GetAuthToken checks the authentication POSTed and issues a token if it's
// valid.
func GetAuthToken(a *AppContext) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		type PostData struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		decoder := json.NewDecoder(r.Body)
		pd := PostData{}
		err := decoder.Decode(&pd)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(FormatError("Unable to decode JSON request."))
			log.Println(err)
			return
		}

		user, err := a.Store.GetUser(pd.Username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(FormatError("error accessing database"))
			log.Println(err)
			return
		}

		if user == nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(FormatError("Unable to find user."))
			log.Println(err)
			return
		}

		passwordMatches, err := user.VerifyPassword(pd.Password)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(FormatError("Unable to authenticate request."))
			log.Println(err)
			return
		}

		if !passwordMatches {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(FormatError("Unable to authenticate request."))
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userid": user.ID,
		})

		tokenString, err := token.SignedString([]byte(a.Config.SigningKey))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(FormatError("Unable to sign token."))
			log.Println(err)
			return
		}

		w.Write([]byte(fmt.Sprintf("{\"token\":\"%s\"}", tokenString)))

	})

	return handler
}
