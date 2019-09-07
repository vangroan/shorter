package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

const urlRegex string = `(https?:\/\/(?:www\.|[^w]{3})[a-zA-Z0-9][a-zA-Z0-9-]+[a-zA-Z0-9]\.[^\s]{2,}|www\.[a-zA-Z0-9][a-zA-Z0-9-]+[a-zA-Z0-9]\.[^\s]{2,}|https?:\/\/(?:www\.|[^w]{3})[a-zA-Z0-9]+\.[^\s]{2,}|www\.[a-zA-Z0-9]+\.[^\s]{2,})`

func validateURI(uri string) error {
	if strings.HasPrefix(uri, "data:") {
		return fmt.Errorf("Data URI not allowed")
	}

	if strings.HasPrefix(uri, "javascript:") {
		return fmt.Errorf("Javascript URI not allowed")
	}

	if !strings.HasPrefix(uri, "http://") && !strings.HasPrefix(uri, "https://") {
		return fmt.Errorf("URI must start with either http or https")
	}

	match, err := regexp.Match(urlRegex, []byte(uri))
	if err != nil {
		return err
	}
	if !match {
		return fmt.Errorf("URI format invalid")
	}

	return nil
}

func respondJSON(w http.ResponseWriter, statusCode int, correlationID string, payload interface{}) error {
	log.
		WithFields(log.Fields{
			"correlation-id": correlationID,
			"status":         statusCode,
			"response":       payload,
		}).
		Trace("Response")

	data, err := json.MarshalIndent(payload, "", "	")
	if err != nil {
		log.Error("Failed to marshal JSON response")
		return err
	}

	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", "application/json")
	_, err = fmt.Fprintln(w, string(data))
	if err != nil {
		log.Error("Failed to write JSON data to response")
		return err
	}

	return nil
}

// Controller is a container for the dependencies
// of the andler functions.
type Controller struct {
	store        Storage
	baseShortURL string
}

// NewController creates a new `Controller`.
//
// Can fail when required dependencies are `nil`.
func NewController(store Storage, baseShortURL string) (*Controller, error) {
	if store == nil {
		return nil, fmt.Errorf("Storage is nil")
	}

	return &Controller{
		store:        store,
		baseShortURL: baseShortURL,
	}, nil
}

// Home is the default page
func (c *Controller) Home(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("./assets/index.html")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 Not Found: %s\n", err.Error())
		respondJSON(w,
			http.StatusInternalServerError,
			r.Header.Get("Correlation-ID"),
			ErrorResponse{
				Message: err.Error(),
			})
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header()["Content-Type"] = []string{mime.TypeByExtension(".html")}
	fmt.Fprintln(w, string(content))
}

// RedirectShort performs a lookup of the
// shortened URL given, and maps it to a long URL.
//
// Responds with a HTTP redirect to the long URL.
func (c *Controller) RedirectShort(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	correlationID := r.Header.Get("Correlation-ID")

	path, ok := vars["path"]
	if !ok {
		log.
			WithFields(log.Fields{
				"path":           path,
				"correlation-id": correlationID,
			}).
			Error("Failed to retrieve path from URL")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Failed to retrieve path from URL")
		return
	}

	log.
		WithFields(log.Fields{
			"path":           path,
			"correlation-id": correlationID,
		}).
		Trace("Redirecting")

	id := DecodeNumber(path)
	location, err := c.store.GetLocation(id)
	if err != nil {
		log.
			WithFields(log.Fields{
				"path":           path,
				"correlation-id": correlationID,
			}).
			Error(err)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "404 Not Found")
		return
	}

	http.Redirect(w, r, location.URL, http.StatusPermanentRedirect)
}

// NotFound returns an HTTP Not Found Error
func (c *Controller) NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(w, "404 Not Found")
}

// CreateURI takes a number and returns a shortened representation
func (c *Controller) CreateURI(w http.ResponseWriter, r *http.Request) {
	correlationID := r.Header.Get("Correlation-ID")
	log.
		WithFields(log.Fields{
			"correlation-id": correlationID,
		}).
		Info("Creating Short URI")

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondJSON(w,
			http.StatusInternalServerError,
			correlationID,
			ErrorResponse{
				Message: fmt.Sprint("Failed to read request body: ", err.Error()),
			})
		return
	}

	var create CreateRequest
	err = json.Unmarshal(content, &create)
	if err != nil {
		respondJSON(w,
			http.StatusInternalServerError,
			correlationID,
			ErrorResponse{
				Message: fmt.Sprint("Failed to read request body: ", err.Error()),
			})
		return
	}

	create.URL = strings.TrimSpace(create.URL)

	// Request must contain a URL
	if create.URL == "" {
		respondJSON(w,
			http.StatusBadRequest,
			correlationID,
			ErrorResponse{
				Message: "No URL provided in POST body",
			})
		return
	}

	// Validate for security
	err = validateURI(create.URL)
	if err != nil {
		respondJSON(w,
			http.StatusBadRequest,
			correlationID,
			ErrorResponse{
				Message: err.Error(),
			})
		return
	}

	location := Location{
		URL:       create.URL,
		TTL:       60 * 60 * 48, // Two days in seconds
		CreatedAt: time.Now(),
	}

	// Save will populate short URL in location
	err = c.store.SaveLocation(&location)
	if err != nil {
		respondJSON(w,
			http.StatusInternalServerError,
			correlationID,
			ErrorResponse{
				Message: fmt.Sprint("Error saving URL: ", err.Error()),
			})
		return
	}

	short := EncodeNumber(location.ID)
	if strings.HasSuffix(c.baseShortURL, "/") {
		short = c.baseShortURL + short
	} else {
		short = c.baseShortURL + "/" + short
	}

	response := CreateResponse{
		Short: short,
		Long:  location.URL,
	}
	respondJSON(w, http.StatusCreated, correlationID, response)
}
