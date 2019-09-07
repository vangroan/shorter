package main

import (
	"time"
)

// CreateRequest is the information
// needed to create a new URL mapping.
//
// The time-to-live is optional, and
// indicates when the mapping expires.
type CreateRequest struct {
	URL string `json:"url"`
	TTL int    `json:"ttl"`
}

// CreateResponse is the information
// that is returned when creating a new
// URL mapping.
type CreateResponse struct {
	Short string `json:"short"`
	Long  string `json:"long"`
}

// Location is the database record
// that stores the target URL of the
// mapping.
//
// The short URL is not stored. When a
// route is handled, it's decoded to
// a location's primary key.
type Location struct {
	ID        uint64    `gorm:"PRIMARY_KEY;AUTO_INCREMENT;Column:id"`
	URL       string    `gorm:"Column:url"`
	TTL       uint64    `gorm:"Column:ttl"`
	CreatedAt time.Time `gorm:"Column:created_at"`
}

// ErrorResponse is the common
// error shape for endpoints that
// are meant to be used as API
// endpoints, and not UI pages.
type ErrorResponse struct {
	Message string `json:"message"`
}
