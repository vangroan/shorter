package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"

	log "github.com/sirupsen/logrus"

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
	sqlite        string
	baseURL       string
	port          string
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
		correlationID := uuid.New()
		ctxWithCorrelation := context.WithValue(r.Context(), CorrelationKey, correlationID)

		logger := log.
			WithFields(log.Fields{
				"correlation-id": correlationID.String(),
				"method":         r.Method,
				"url":            r.RequestURI,
				"operation":      "Request",
			})
		logger.Info("Incoming Request")

		ctxWithLogger := context.WithValue(ctxWithCorrelation, LogKey, logger)
		req := r.WithContext(ctxWithLogger)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, req)
	})
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.TraceLevel)

	log.Println("Starting")

	// Load config
	config := getConfig()

	log.Printf("Base URL: %s", config.baseURL)

	// Opening database
	db, err := gorm.Open("sqlite3", config.sqlite)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.AutoMigrate(
		&Location{},
		&UserAccount{},
		&Subscription{},
	).Error
	if err != nil {
		panic(err)
	}

	db.BlockGlobalUpdate(true)

	// Setup storage service
	store := NewDBStorage(db)

	// Start background jobs
	cancel := make(chan struct{})
	var wg sync.WaitGroup

	TimeToLiveJob(cancel, &wg, 10*time.Second, &store)

	// Setup controller
	ctrl, err := NewController(&store, config.baseURL)
	if err != nil {
		panic(err)
	}

	// Setup Router
	r := mux.NewRouter()

	// This will serve files under http://localhost:8000/assets/<filename>
	r.PathPrefix("/assets").
		Methods("GET").
		Handler(http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))
	r.Path("/").Methods("GET").HandlerFunc(ctrl.Home)
	r.Path("/{path}").Methods("GET").HandlerFunc(ctrl.RedirectShort)
	r.Path("/").Methods("POST").Handler(AuthMiddleware(http.HandlerFunc(ctrl.CreateURI)))
	r.Use(loggingMiddleware)
	r.NotFoundHandler = loggingMiddleware(http.HandlerFunc(ctrl.NotFound))
	http.Handle("/", r)

	serve := http.Server{
		Handler:      r,
		Addr:         ":" + config.port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Catch Ctrl-C for graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(cancel)
		defer close(c)

		for sig := range c {
			// sig is a ^C, handle it
			log.Println("Received signal:", sig.String())
			break
		}

		ctx, cncl := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
		defer cncl()

		log.Println("Server shutting down")
		if err := serve.Shutdown(ctx); err != nil {
			log.Fatalf("Shutdown(): %s", err)
		}
	}()

	// Start server
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := serve.ListenAndServe(); err != nil {
			log.Printf("ListenAndServe(): %s", err)
		}
		log.Println("Server down")
	}()

	wg.Wait()
	log.Println("Done")
}
