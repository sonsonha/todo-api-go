package http

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sonsonha/todo-api-go/internal/store"
)

type Handler struct {
	DB *sql.DB
	Q  *store.Queries
}

// NewHandler creates a new handler with DB connection
func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		DB: db,
		Q:  store.New(db),
	}
}

// RegisterRoutes adds routes to chi router
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/todos", h.CreateTodo)
	r.Get("/todos", h.ListTodos)
	r.Get("/todos/{id}", h.GetTodo)
	r.Patch("/todos/{id}", h.UpdateTodo)
	r.Delete("/todos/{id}", h.DeleteTodo)
}

// CreateTodo handler (POST /todos)
func (h *Handler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	todo, err := h.Q.CreateTodo(context.Background(), input.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// ListTodos handler (GET /todos?limit=10&offset=0)
func (h *Handler) ListTodos(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, _ := strconv.Atoi(limitStr)
	if limit == 0 {
		limit = 10
	}
	offset, _ := strconv.Atoi(offsetStr)

	todos, err := h.Q.ListTodos(context.Background(), store.ListTodosParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

// GetTodo handler (GET /todos/{id})
func (h *Handler) GetTodo(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	todo, err := h.Q.GetTodo(context.Background(), int64(id))
	if err != nil {
		http.Error(w, "todo not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// UpdateTodo handler (PATCH /todos/{id})
func (h *Handler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var input struct {
		Title  string `json:"title"`
		IsDone bool   `json:"is_done"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	todo, err := h.Q.UpdateTodo(context.Background(), store.UpdateTodoParams{
		ID:     int64(id),
		Title:  input.Title,
		IsDone: input.IsDone,
	})
	if err != nil {
		http.Error(w, "todo not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// DeleteTodo handler (DELETE /todos/{id})
func (h *Handler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = h.Q.DeleteTodo(context.Background(), int64(id))
	if err != nil {
		http.Error(w, "failed to delete todo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
