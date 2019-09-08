package main

// contextKey is the package wide type
// for value keys in the request context.
type contextKey int

const (
	// AuthUserKey is the key for the current
	// request scoped authenticated user.
	AuthUserKey contextKey = iota

	// LogKey is the key for the current logging
	// context, which contains request scoped
	// fields, such as correlation ID.
	LogKey contextKey = iota

	// CorrelationKey is the key for the current
	// request's correlation ID.
	CorrelationKey contextKey = iota
)
