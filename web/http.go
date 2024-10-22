package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// encode writes the provided data to the response writer.
func encode[T any](w http.ResponseWriter, status int, data T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return fmt.Errorf("encode: %w", err)
	}
	return nil
}

// decode reads the request body and decodes it into the provided value.
func decode[T any](r *http.Request) (T, error) {
	var data T
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return data, fmt.Errorf("decode: %w", err)
	}

	return data, nil
}
