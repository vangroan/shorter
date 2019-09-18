package main

import (
	"time"

	"github.com/jinzhu/gorm"
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
//
// A location with no subscription
// implicitly belongs to an anonymous
// user.
type Location struct {
	ID           uint64 `gorm:"PRIMARY_KEY;AUTO_INCREMENT"`
	URL          string
	TTL          uint64
	Subscription *Subscription
	CreatedAt    time.Time
}

// ErrorResponse is the common
// error shape for endpoints that
// are meant to be used as API
// endpoints, and not UI pages.
type ErrorResponse struct {
	Message string `json:"message"`
}

// UserAccount represents a user's profile
// with the service.
type UserAccount struct {
	gorm.Model
	Subscriptions []*Subscription `gorm:"many2many:account_subscriptions;"`
	IsAdmin       bool
	isAnonymous   bool `gorm:"-"` // ignore this field
}

// Subscription keeps user accounts seperate
// from their URLs, so a subscription can be
// created and deleted without affecting the account
// itself.
//
// Multiple accounts can be added to multiple subscriptions,
// to allow for team access.
//
// Limits and restrictions are based on subscription.
type Subscription struct {
	gorm.Model
	Accounts  []*UserAccount `gorm:"many2many:account_subscriptions;"`
	Locations []Location
}

// IsAnonymous indicates whether the user
// is authenticated or not.
func (u *UserAccount) IsAnonymous() bool {
	return u.isAnonymous
}

// NewAnonymousUser creates a special user
// who is not identified or authenticated.
func NewAnonymousUser() *UserAccount {
	return &UserAccount{
		Model: gorm.Model{
			ID:        0,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			DeletedAt: nil,
		},
		Subscriptions: nil,
		IsAdmin:       false,
		isAnonymous:   true,
	}
}
