package auth

import (
	"net/http"
)

// MockJWTAuthMiddleware is a mock middleware for JWT authentication
func MockJWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Bypass actual JWT validation for testing
		next.ServeHTTP(w, r)
	})
}
