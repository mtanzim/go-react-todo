package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/rs/cors"
	_ "modernc.org/sqlite"
)

type NewTodo struct {
	Name string `json:"name"`
}

type Todo struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}

func JSONError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func main() {

	db, err := sql.Open("sqlite", "./todo.db")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS todos (id INTEGER PRIMARY KEY, name TEXT, completed BOOLEAN)")
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world"))
	})
	r.Post("/api/v1/todo/", func(w http.ResponseWriter, r *http.Request) {
		var todo Todo
		err := json.NewDecoder(r.Body).Decode(&todo)
		if err != nil {
			JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		res, err := db.Exec("INSERT INTO todos (name, completed) VALUES (?, ?)", todo.Name, false)
		if err != nil {
			JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		id, err := res.LastInsertId()
		if err != nil {
			JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		insertedTodo := Todo{
			ID:        id,
			Name:      todo.Name,
			Completed: todo.Completed,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(insertedTodo)

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
		row := db.QueryRow("SELECT id, name, completed FROM todos WHERE id = ?", id)
		var todo Todo
		err := row.Scan(&todo.ID, &todo.Name, &todo.Completed)
		if err != nil {
			JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(todo)
	})
	r.Get("/api/v1/todo/", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name, completed FROM todos")
		if err != nil {
			JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var todos []Todo
		defer rows.Close()
		for rows.Next() {
			var todo Todo
			err := rows.Scan(&todo.ID, &todo.Name, &todo.Completed)
			if err != nil {
				JSONError(w, err.Error(), http.StatusInternalServerError)
				return
			}
			todos = append(todos, todo)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(todos)
	})
	http.ListenAndServe(":8080", cors.Default().Handler(r))
}
