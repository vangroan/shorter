package main

import (
	"time"
)

type CreateRequest struct {
	URL string `json:"url"`
	TTL int    `json:"ttl`
}

type CreateResponse struct {
	Short string `json:"short"`
	Long  string `json:"long"`
}

type Location struct {
	ID        uint64 `gorm:"PRIMARY_KEY;AUTO_INCREMENT"`
	URL       string
	CreatedAt time.Time
}

type ErrorResponse struct {
	Message string `json:"message"`
}
