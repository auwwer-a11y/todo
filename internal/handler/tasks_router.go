package handler

import (
	"github.com/go-chi/chi/v5"
	"log/slog"
	"github.com/auwwer-a11y/todo/internal/service"
	"net/http"
	"context"
	"encoding/json"
)

type TasksRouter struct {
	taskService *service.TaskService
	noteService *service.NoteService
	logger      *slog.Logger
	authServiceURL string
}

func NewTasksRouter(taskService *service.TaskService, noteService *service.NoteService, logger *slog.Logger, authServiceURL string) http.Handler {
	h := &TasksRouter{
		taskService: taskService,
		noteService: noteService,
		logger: logger,
		authServiceURL: authServiceURL, 
	}

	r := chi.NewRouter()

	r.Get("/healthz", h.handleHealthz)
	r.Get("/readyz", h.handleReadyz)

	r.Group(func(r chi.Router) {
		r.Use(h.authMiddleware)

		r.Get("/api/tasks/{id}", h.handleGetTask)

		r.Get("/api/tasks", h.handleGetTasks)
		r.Post("/api/tasks", h.handleCreateTask)
		r.Put("/api/tasks/{id}", h.handleUpdateTask)
		r.Delete("/api/tasks/{id}", h.handleDeleteTask)

		r.Get("/api/tasks/{id}/notes", h.handleGetNotes)
		r.Post("/api/tasks/{id}/notes", h.handleCreateNote)
		r.Delete("/api/notes/{id}", h.handleDeleteNote)
	})
	return r
}

func (h *TasksRouter) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		req, _ := http.NewRequestWithContext(r.Context(), "POST", h.authServiceURL+"/api/validate", nil)
		req.Header.Set("Authorization", authHeader)

		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var result map[string]string
		json.NewDecoder(resp.Body).Decode(&result)
		userID := result["user_id"]
		ctx := context.WithValue(r.Context(), "userID", userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
