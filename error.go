package main

import (
	"encoding/json"
)

// ErrorMessage represents an error in the format that will be returned on API
// calls that result in an error.
type ErrorMessage struct {
	Message string `json:"message"`
}

// FormatError formats a simple error message into a JSON resoponse.
func FormatError(errorMsg string) []byte {
	j, err := json.Marshal(ErrorMessage{Message: errorMsg})

	if err != nil {
		j = []byte("{\"message\":\"An error occurred.\"}")
	}

	return j
}
