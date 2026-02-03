package handler

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type ctxKey string

const RequestIDKey ctxKey = "request_id"

func GetRequestID(ctx context.Context) string {
	if v := ctx.Value(RequestIDKey); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}

	return ""
}

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := r.Header.Get("X-Request-ID")
		if rid == "" {
			rid = uuid.NewString()
		}

		// Put into context
		ctx := context.WithValue(r.Context(), RequestIDKey, rid)
		r = r.WithContext(ctx)

		// Return header back (useful for tracing)
		w.Header().Set("X-Request-Id", rid)

		next.ServeHTTP(w, r)
	})
}
