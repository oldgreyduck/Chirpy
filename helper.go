package main

import (
	"net/http"
	"encoding/json"
)


func respondWithError(w http.ResponseWriter, code int, msg string) error {
	//...
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(b)
	return err
}

func cleanChirp(body string) string {
	//...
}
