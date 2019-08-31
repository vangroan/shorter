package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
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

	// Opening database
	db, err := gorm.Open("sqlite3", "data.sqlite3")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err = db.AutoMigrate(&Location{}).Error; err != nil {
		panic(err)
	}

	// Setup storage service
	store := NewSQLitestorage(db)

	// Setup controller
	ctrl, err := NewController(&store)
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
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	serve.ListenAndServe()

	log.Println("Stopping")
}
