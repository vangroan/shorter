package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func respondJSON(w http.ResponseWriter, statusCode int, payload interface{}) error {
	data, err := json.MarshalIndent(payload, "", "	")
	if err != nil {
		return err
	}

	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", "application/json")
	_, err = fmt.Fprintln(w, string(data))
	if err != nil {
		return err
	}

	return nil
}

type Controller struct {
	store Storage
}

func NewController(store Storage) (*Controller, error) {
	if store == nil {
		return nil, fmt.Errorf("Storage is nil")
	}

	return &Controller{
		store: store,
	}, nil
}

// Home is the default page
func (c *Controller) Home(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Hello world")
}

// RedirectShort performs a lookup of the
// shortened URL given, and maps it to a long URL.
//
// Responds with a HTTP redirect to the long URL.
func (c *Controller) RedirectShort(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	_, ok := vars["path"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Failed to retrieve path from URL")
		return
	}

	// http.Redirect(w, r, , http.StatusPermanentRedirect)
	// w.WriteHeader(http.StatusOK)
	// fmt.Fprintln(w, DecodeNumber(path))
}

// NotFound returns an HTTP Not Found Error
func (c *Controller) NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(w, "404 Not Found")
}

// CreateURI takes a number and returns a shortened representation
func (c *Controller) CreateURI(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Message: fmt.Sprint("Failed to read request body: ", err.Error()),
		})
		return
	}

	var create CreateRequest
	err = json.Unmarshal(content, &create)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Message: fmt.Sprint("Failed to read request body: ", err.Error()),
		})
		return
	}

	// Request must contain a URL
	if create.URL == "" {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "No URL provided in POST body",
		})
		return
	}

	location := Location{
		URL:       create.URL,
		CreatedAt: time.Now(),
	}

	// Save will populate short URL in location
	err = c.store.SaveLocation(&location)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Message: fmt.Sprint("Error saving URL: ", err.Error()),
		})
		return
	}

	response := CreateResponse{
		Short: EncodeNumber(location.ID),
		Long:  location.URL,
	}
	respondJSON(w, http.StatusCreated, response)
}
