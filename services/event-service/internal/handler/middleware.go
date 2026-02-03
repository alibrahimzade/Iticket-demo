package handler

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)
		rid := GetRequestID(r.Context())
		log.Printf("%s %s %s %s", rid, r.Method, r.URL.Path, time.Since(start))
	})
}

// RecoveryMiddleware prevents server crash on panic
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("panic recovered:", err)
				respondError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "unexpected server error", nil)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
