package middlewares

import (
	"log"
	"net/http"
)

// Use Decorators for the middlewares.

// SetLogger logs the request.
func SetLogger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("\n%s %s%s %s", r.Method, r.Host, r.RequestURI, r.Proto)
		next(w, r)
	}
}

// SetJSON sets the content-type as application/json.
func SetJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}
