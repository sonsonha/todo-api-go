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
