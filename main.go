package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func main() {
	log.Println("Starting")

	r := mux.NewRouter()
	r.Path("/").Methods("GET").HandlerFunc(HomeHandler)
	r.Path("/").Methods("POST").HandlerFunc(CreateURIHandler)
	r.Use(loggingMiddleware)
	r.NotFoundHandler = loggingMiddleware(http.HandlerFunc(NotFoundHandler))
	http.Handle("/", r)

	serve := http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	serve.ListenAndServe()

	log.Println("Stopping")
}
