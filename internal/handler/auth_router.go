package handler

import (
	"github.com/go-chi/chi/v5"
	"log/slog"
	"github.com/auwwer-a11y/todo/internal/service"
	"net/http"
	"context"
	"strings"
	"github.com/auwwer-a11y/todo/pkg/jwt"
	"encoding/json"
)

type AuthRouter struct {
	userService *service.UserService
	logger      *slog.Logger
	jwtSecret   string
}

func NewAuthRouter(userService *service.UserService, logger *slog.Logger, jwtSecret string) http.Handler {
	h := &AuthRouter{
		userService: userService,
		logger: logger,
		jwtSecret: jwtSecret,
	}

	r := chi.NewRouter()

	r.Post("/api/register", h.handleRegister)
	r.Post("/api/login", h.handleLogin)
	r.Post("/api/validate", h.handleValidate)

	r.Get("/healthz", h.handleHealthz)
	r.Get("/readyz", h.handleReadyz)

	return r
}

func (h *AuthRouter) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r * http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		userID, err := jwt.Validate(tokenString, h.jwtSecret)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *AuthRouter) handleValidate(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	userID, err := jwt.Validate(tokenString, h.jwtSecret)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"user_id": userID})
}