package main

import "net/http"

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthz)
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
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