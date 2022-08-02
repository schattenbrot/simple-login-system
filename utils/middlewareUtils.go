package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func MiddlewareBodyDecoder(r *http.Request, obj interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.NewDecoder(bytes.NewBuffer(body)).Decode(obj)
	if err != nil {
		return err
	}
	r.Body.Close()
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	return nil
}
