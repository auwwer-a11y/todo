package handler

import (
	"encoding/json"
	"net/http"
	"errors"
	"github.com/auwwer-a11y/todo/internal/domain"
)

func (h *Router) handleRegister(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var req struct {
		Name string `json:"name"`
		Email string `json:"email"`
		Password string `json:"password"`
	}

	if err := decoder.Decode(&req); err != nil {
		h.logger.Error("Failed to decode register request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userService.Register(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}
		h.logger.Error("Failed to register user", "error", err)
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	resp := struct {
		ID string `json:"id"`
	}{
		ID: user.ID,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *Router) handleLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var req struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	if err := decoder.Decode(&req); err != nil {
		h.logger.Error("Failed to decode login request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := h.userService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		h.logger.Error("Failed to login user", "error", err)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	resp := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
