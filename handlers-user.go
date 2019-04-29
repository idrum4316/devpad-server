package main

import (
	"encoding/json"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/idrum4316/devpad-server/internal/user"
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
			_, _ = w.Write(FormatError("Unable to decode JSON request."))
			log.Println(err)
			return
		}

		u, err := a.Store.GetUser(pd.Username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(FormatError("error accessing database"))
			log.Println(err)
			return
		}

		if u == nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(FormatError("Unable to find user."))
			return
		}

		passwordMatches, err := u.VerifyPassword(pd.Password)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(FormatError("Unable to authenticate request."))
			log.Println(err)
			return
		}

		if !passwordMatches {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(FormatError("Unable to authenticate request."))
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userid": u.ID,
		})

		tokenString, err := token.SignedString([]byte(a.Config.SigningKey))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(FormatError("Unable to sign token."))
			log.Println(err)
			return
		}

		response := map[string]interface{}{
			"token":    tokenString,
			"is_admin": u.Admin,
			"username": u.ID,
		}

		responseJSON, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(FormatError("unable to encode response"))
			log.Println(err)
			return
		}

		_, _ = w.Write(responseJSON)

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
			_, _ = w.Write(FormatError("error parsing authentication token"))
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
			_, _ = w.Write(FormatError("Unable to decode JSON request."))
			log.Println(err)
			return
		}

		// Make sure both new passwords match
		if pd.New != pd.Confirm {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(FormatError("passwords do not match"))
			return
		}

		// Fetch the user account from the datastore
		u, err := a.Store.GetUser(userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(FormatError("error accessing database"))
			log.Println(err)
			return
		}

		// Verify that the current password is valid
		currentPasswordMatches, err := u.VerifyPassword(pd.Current)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(FormatError("could not verify current password"))
			return
		}
		if !currentPasswordMatches {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(FormatError("current password is not valid"))
			return
		}

		// Set the new password
		err = u.SetPassword(pd.New)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(FormatError("unable to set new password"))
			log.Println(err)
			return
		}

		// Update the user in the data store
		err = a.Store.UpdateUser(u)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(FormatError("unable to set new password"))
			log.Println(err)
			return
		}

		// Return a success message
		return

	})

	return RequireAuth(handler, a)
}

// CreateUserHandler creates a new user in the data store. The requesting user
// must be an admin.
func CreateUserHandler(a *AppContext) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Parse the user id from the request's auth token
		userID, err := a.GetUserIDFromRequest(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(FormatError("error parsing authentication token"))
			log.Println(err)
			return
		}

		// Fetch the user account from the datastore
		u, err := a.Store.GetUser(userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(FormatError("error accessing database"))
			log.Println(err)
			return
		}

		// If the user is not an admin, they aren't authorized to create
		// accounts.
		if !u.Admin {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Parse the body of the POST request
		type PostData struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Confirm  string `json:"confirm"`
			Admin    bool   `json:"is_admin"`
		}
		decoder := json.NewDecoder(r.Body)
		pd := PostData{}
		err = decoder.Decode(&pd)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(FormatError("Unable to decode JSON request."))
			log.Println(err)
			return
		}

		// Make sure both new passwords match
		if pd.Password != pd.Confirm {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(FormatError("passwords do not match"))
			return
		}

		// Create the new user account
		newUser := user.User{
			ID:    pd.Username,
			Admin: pd.Admin,
		}

		// Set the new password
		err = newUser.SetPassword(pd.Password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(FormatError("unable to set new password"))
			log.Println(err)
			return
		}

		// Update the user in the data store
		err = a.Store.CreateUser(&newUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(FormatError("unable to set new password"))
			log.Println(err)
			return
		}

		return

	})

	return RequireAuth(handler, a)
}
