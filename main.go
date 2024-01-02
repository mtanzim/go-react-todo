package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

func main() {

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world"))
	})
	r.Post("/api/v1/todo/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Create todo"))
	})
	r.Put("/api/v1/todo/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		response := "Update todo with ID: " + id
		w.Write([]byte(response))
	})
	r.Delete("/api/v1/todo/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		response := "delete todo with ID: " + id
		w.Write([]byte(response))
	})
	r.Delete("/api/v1/todo/", func(w http.ResponseWriter, r *http.Request) {
		response := "delete all todo"
		w.Write([]byte(response))
	})
	r.Get("/api/v1/todo/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		response := "get todo with ID: " + id
		w.Write([]byte(response))
	})
	r.Get("/api/v1/todo/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Get all todos"))
	})
	http.ListenAndServe(":8080", r)
}
