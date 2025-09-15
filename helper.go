package main

import (
	"net/http"
	"encoding/json"
	"strings"
)

var bannedWords = map[string]struct{}{
    "kerfuffle": {},
    "sharbert":  {},
    "fornax":    {},
}


func respondWithError(w http.ResponseWriter, code int, msg string) error {
    b, err := json.Marshal(map[string]string{"error": msg})
    if err != nil {
        return err
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    _, err = w.Write(b)
    return err
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
    parts := strings.Split(body, " ")
    for i, p := range parts {
        if _, ok := bannedWords[strings.ToLower(p)]; ok {
            parts[i] = "****"
        }
    }
    return strings.Join(parts, " ")
}
