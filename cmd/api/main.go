package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	apihttp "github.com/sonsonha/todo-api-go/internal/http"
	"github.com/sonsonha/todo-api-go/internal/store"
)

func main() {
	// Connect DB
	db, err := store.Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	r := chi.NewRouter()

	// Register health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Register todo routes
	h := apihttp.NewHandler(db)
	h.RegisterRoutes(r)

	fmt.Println("🚀 Server running on :8080")
	http.ListenAndServe(":8080", r)
}
