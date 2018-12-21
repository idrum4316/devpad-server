package main

import (
	"encoding/json"
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

		response := map[string]interface{}{
			"token":    tokenString,
			"is_admin": user.Admin,
			"username": user.ID,
		}

		responseJSON, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(FormatError("unable to encode response"))
			log.Println(err)
			return
		}

		w.Write(responseJSON)

	})

	return handler
}

// ChangePasswordHandler updates the user's password
func ChangePasswordHandler(a *AppContext) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Parse the user id from the request's auth token
		userID, err := a.GetUserIDFromRequest(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(FormatError("error parsing authentication token"))
			log.Println(err)
			return
		}

		// Parse the body of the POST request
		type PostData struct {
			Current string `json:"current"`
			New     string `json:"new"`
			Confirm string `json:"confirm"`
		}
		decoder := json.NewDecoder(r.Body)
		pd := PostData{}
		err = decoder.Decode(&pd)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(FormatError("Unable to decode JSON request."))
			log.Println(err)
			return
		}

		// Make sure both new passwords match
		if pd.New != pd.Confirm {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(FormatError("passwords do not match"))
			return
		}

		// Fetch the user account from the datastore
		user, err := a.Store.GetUser(userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(FormatError("error accessing database"))
			log.Println(err)
			return
		}

		// Verify that the current password is valid
		currentPasswordMatches, err := user.VerifyPassword(pd.Current)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(FormatError("could not verify current password"))
			return
		}
		if !currentPasswordMatches {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(FormatError("current password is not valid"))
			return
		}

		// Set the new password
		err = user.SetPassword(pd.New)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(FormatError("unable to set new password"))
			log.Println(err)
			return
		}

		// Update the user in the data store
		err = a.Store.UpdateUser(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(FormatError("unable to set new password"))
			log.Println(err)
			return
		}

		// Return a success message
		return

	})

	return RequireAuth(handler, a)
}
