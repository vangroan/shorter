package main

import (
	"context"
	"net/http"
)

// AuthMiddleware is a middleware that extracts
// an authenticated user's identity from an
// incoming request, queries the user store, and
// adds the user to the request.
//
// If the JWT in the header is invalid, an
// error response will be returned.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Validate Token
		// TODO: Get user from DB
		ctx := context.WithValue(r.Context(), AuthUserKey, NewAnonymousUser())
		req := r.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}
