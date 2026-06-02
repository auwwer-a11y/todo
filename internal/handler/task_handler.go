package handler

import (
	"encoding/json"
	"net/http"
	"time"
	"errors"
	"github.com/auwwer-a11y/todo/internal/domain"
	"github.com/go-chi/chi/v5"
)

func (h *Router) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var req struct {
		Title string `json:"title"`
		Description string `json:"description"`
		Deadline *time.Time `json:"deadline"`
	}

	if err := decoder.Decode(&req); err != nil {
		h.logger.Error("Failed to decode create task request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(string)

	task, err := h.taskService.Create(r.Context(), userID, req.Title, req.Description, req.Deadline)
	if err != nil {
		h.logger.Error("Failed to create task", "error", err)
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	resp := struct {
		ID string `json:"id"`
	}{
		ID: task.ID,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *Router) handleGetTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID := r.Context().Value("userID").(string)

	task, err := h.taskService.GetByID(r.Context(), userID, id)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, domain.ErrTaskForbidden) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
		h.logger.Error("Failed to get task", "error", err)
		http.Error(w, "Failed to get task", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func (h *Router) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	decoder := json.NewDecoder(r.Body)

	var req struct {
		Title string `json:"title"`
		Description string `json:"description"`
		Status domain.TaskStatus `json:"status"`
		Deadline *time.Time `json:"deadline"`
	}

	if err := decoder.Decode(&req); err != nil {
		h.logger.Error("Failed to decode update task request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(string)

	task, err := h.taskService.Update(r.Context(), userID, id, req.Title, req.Description, req.Status, req.Deadline)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, domain.ErrTaskForbidden) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		h.logger.Error("Failed to update task", "error", err)
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func (h *Router) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	id := chi.URLParam(r, "id")

	err := h.taskService.Delete(r.Context(), userID, id)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, domain.ErrTaskForbidden) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		h.logger.Error("Failed to delete task", "error", err)
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}