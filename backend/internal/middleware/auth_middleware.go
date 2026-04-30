package middleware

import (
	"context"
	"net/http"
	"strings"

	"product-management/backend/internal/auth"
	"product-management/backend/internal/response"
)

type contextKey string

const UserContextKey contextKey = "auth_user"

func Auth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "missing authorization header", nil)
				return
			}

			parts := strings.SplitN(header, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid authorization format", nil)
				return
			}

			claims, err := auth.ParseToken(secret, parts[1])
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid token", nil)
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
