package utils

import (
	"encoding/json"
	"net/http"
)

// writeJSON is the helper function for sending back an HTTP response.
func WriteJSON(w http.ResponseWriter, status int, data ...interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if len(data) == 0 {
		w.Write([]byte(""))
		return nil
	}

	js, err := json.Marshal(data[0])
	if err != nil {
		return err
	}

	w.Write(js)

	return nil
}

// errorJSON is the helper function for creating an error message.
// This internally then runs the writeJSON function to send the HTTP response.
func ErrorJSON(w http.ResponseWriter, err error, status ...int) {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	type jsonError struct {
		StatusCode int    `json:"statusCode"`
		Message    string `json:"message"`
	}

	theError := jsonError{
		StatusCode: statusCode,
		Message:    err.Error(),
	}

	WriteJSON(w, statusCode, theError)
}
