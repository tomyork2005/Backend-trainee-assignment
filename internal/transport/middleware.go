package transport

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
)

const (
	bearerPrefix        = "Bearer "
	authorizationHeader = "Authorization"
	userIDContextKey    = "UserID"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.With(
			slog.String("METHOD", r.Method),
			slog.String("URL", r.URL.String())).Info("handle logger")
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get(authorizationHeader)
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			slog.Error("Invalid authorization header", slog.String("AuthorizationHeader", authHeader))

			resp := ErrorResponse{
				Errors: "Invalid authorization header",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(resp)
			return
		}

		token := strings.TrimPrefix(authHeader, bearerPrefix)

		userID, err := h.authService.ParseToken(r.Context(), token)
		if err != nil {
			slog.Error("Invalid token", slog.String("token", token), slog.String("error", err.Error()))

			resp := ErrorResponse{
				Errors: "Invalid token",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(resp)
			return
		}

		ctx := context.WithValue(r.Context(), userIDContextKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
