package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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

// CreateURIHandler takes a number and returns a shortened representation
func CreateURIHandler(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to read request body:", err.Error())
		return
	}

	num, err := strconv.Atoi(string(content))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Failed to convert request body to number:", err.Error())
		return
	}

	shortened := EncodeNumber(uint64(num))
	log.Println("Encoded", uint64(num), "to", shortened)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, shortened)
}
