package handler

import (
	"github.com/go-chi/chi/v5"
	"log/slog"
	"github.com/auwwer-a11y/todo/internal/service"
	"net/http"
)

type Router struct {
	userService *service.UserService
	taskService *service.TaskService
	noteService *service.NoteService
	logger      *slog.Logger
}

func NewRouter(userService *service.UserService, taskService *service.TaskService, noteService *service.NoteService, logger *slog.Logger) http.Handler {
	h := &Router{
		userService: userService,
		taskService: taskService,
		noteService: noteService,
		logger: logger,
	}

	r := chi.NewRouter()

	r.Post("/api/register", h.handleRegister)
	r.Post("/api/login", h.handleLogin)

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

