package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const (
	defaultSQLite  string = "data/db.sqlite3"
	defaultBaseURL string = "http://localhost"
	defaultPort    string = "8000"
)

type config struct {
	sqlite  string
	baseURL string
	port    string
}

func getConfig() config {
	sqlite := os.Getenv("SHORTER_SQLITE")
	if sqlite == "" {
		sqlite = defaultSQLite
	}

	port := os.Getenv("SHORTER_PORT")
	if port == "" {
		port = defaultPort
	}

	baseURL := os.Getenv("SHORTER_BASEURL")
	if baseURL == "" {
		baseURL = defaultBaseURL + ":" + port
	}

	// Ensure folders exist
	os.MkdirAll(filepath.Dir(sqlite), os.ModePerm)

	return config{
		sqlite:  sqlite,
		baseURL: baseURL,
		port:    port,
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func main() {
	log.Println("Starting")

	// Load config
	config := getConfig()

	// Opening database
	db, err := gorm.Open("sqlite3", config.sqlite)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err = db.AutoMigrate(&Location{}).Error; err != nil {
		panic(err)
	}

	// Setup storage service
	store := NewDBStorage(db)

	// Setup controller
	ctrl, err := NewController(&store, config.baseURL)
	if err != nil {
		panic(err)
	}

	// Setup Router
	r := mux.NewRouter()
	r.Path("/").Methods("GET").HandlerFunc(ctrl.Home)
	r.Path("/{path}").Methods("GET").HandlerFunc(ctrl.RedirectShort)
	r.Path("/").Methods("POST").HandlerFunc(ctrl.CreateURI)
	r.Use(loggingMiddleware)
	r.NotFoundHandler = loggingMiddleware(http.HandlerFunc(ctrl.NotFound))
	http.Handle("/", r)

	serve := http.Server{
		Handler:      r,
		Addr:         ":" + config.port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	err = serve.ListenAndServe()
	if err != nil {
		panic(err)
	}

	log.Println("Stopping")
}
