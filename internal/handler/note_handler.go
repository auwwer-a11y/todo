package handler

import (
	"encoding/json"
	"net/http"
	"errors"
	"github.com/auwwer-a11y/todo/internal/domain"
	"github.com/go-chi/chi/v5"
)

func (h *Router) handleCreateNote(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var req struct {
		Text string `json:"text"`
	}

	if err := decoder.Decode(&req); err != nil {
		h.logger.Error("Failed to decode create note request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	taskID := chi.URLParam(r, "id")
	userID := r.Context().Value("userID").(string)

	note, err := h.noteService.Create(r.Context(), taskID, req.Text, userID, nil)
	if err != nil {
		h.logger.Error("Failed to create note", "error", err)
		http.Error(w, "Failed to create note", http.StatusInternalServerError)
		return
	}

	resp := struct {
		ID string `json:"id"`
	}{
		ID: note.ID,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *Router) handleGetNotes(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")

	notes, err := h.noteService.GetAllByTaskID(r.Context(), taskID)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, domain.ErrTaskForbidden) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
		}
		h.logger.Error("Failed to get notes", "error", err)
		http.Error(w, "Failed to get notes", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(notes)
}

func (h *Router) handleDeleteNote(w http.ResponseWriter, r *http.Request) {
	noteID := chi.URLParam(r, "id")

	err := h.noteService.Delete(r.Context(), noteID)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, domain.ErrTaskForbidden) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
		}
		h.logger.Error("Failed to get notes", "error", err)
		http.Error(w, "Failed to get notes", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}