package main

import (
	"net/http"
	"sync/atomic"
	"fmt"
	"encoding/json"
)


type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	cfg := &apiConfig{}
	mux := http.NewServeMux()
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", cfg.middlewareMetricsInc(handler))
	mux.HandleFunc("GET /api/healthz", healthz)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)
	mux.HandleFunc("POST /api/validate_chirp", validate_chirp)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	srv.ListenAndServe()
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
    hits := cfg.fileserverHits.Load()
	html := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", hits)
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(html))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
    cfg.fileserverHits.Store(0)
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}

func validate_chirp(w http.ResponseWriter, r *http.Request) {
    type parameters struct {
        Body string `json:"body"`
    }

    var params parameters
    if err := json.NewDecoder(r.Body).Decode(&params); err != nil || params.Body == "" {
        respondWithError(w, http.StatusBadRequest, "Something went wrong")
        return
    }

    if len(params.Body) > 140 {
        respondWithError(w, http.StatusBadRequest, "Chirp is too long")
        return
    }

    cleaned := cleanChirp(params.Body)
    _ = respondWithJSON(w, http.StatusOK, map[string]string{
        "cleaned_body": cleaned,
    })
}
