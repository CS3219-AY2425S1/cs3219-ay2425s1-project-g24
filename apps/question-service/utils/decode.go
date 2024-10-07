package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// DecodeJSONBody decodes the request body into the given destination
func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&dst)
	if err != nil {
		return fmt.Errorf("Invalid request payload: " + err.Error())
	}

	return nil
}
