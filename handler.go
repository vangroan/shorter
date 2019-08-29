package main

import (
	"fmt"
	"net/http"
)

// HomeHandler is the default page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Hello world")
}

// NotFoundHandler returns an HTTP Not Found Error
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(w, "404 Not Found")
}
