package main

// ContextKey is the package wide type
// for value keys in the request context.
type ContextKey int

const (
	// AuthUserKey is the key for the current
	// request scoped authenticated user.
	AuthUserKey ContextKey = iota

	// LogKey is the key for the current logging
	// context, which contains request scoped
	// fields, such as correlation ID.
	LogKey ContextKey = iota

	// CorrelationKey is the key for the current
	// request's correlation ID.
	CorrelationKey ContextKey = iota
)
