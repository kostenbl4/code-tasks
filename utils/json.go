package utils

import (
	"encoding/json"
	"net/http"
)

func ReadJSON(r *http.Request, data any) error {
	decoder := json.NewDecoder(r.Body)
	// decoder.DisallowUnknownFields() // для пропуска данных только с правильными полями
	return decoder.Decode(data)
}

func WriteJSON(w http.ResponseWriter, data any, status int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, status int, errMsg string) error {
	type err struct {
		Err string `json:"error"`
	}
	return WriteJSON(w, err{errMsg}, status)
}
