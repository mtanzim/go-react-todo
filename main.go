package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
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

type UpdatedTodo struct {
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
	tds := NewTodoService(db)
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS todos (id INTEGER PRIMARY KEY, name TEXT, completed BOOLEAN)")
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		todos, err := tds.GetAllTodos()
		if err != nil {
			errDiv().Render(context.Background(), w)
			return
		}
		todoDiv(todos).Render(context.Background(), w)
	})

	r.Post("/todo", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		newTodoTitle := r.Form.Get("newTodo")
		newTodo, err := tds.AddTodo(newTodoTitle)
		if err != nil {
			errDiv().Render(context.Background(), w)
			return
		}
		todoCard(newTodo).Render(context.Background(), w)

	})

	r.Post("/api/v1/todo", func(w http.ResponseWriter, r *http.Request) {
		var todo Todo
		err := json.NewDecoder(r.Body).Decode(&todo)
		if err != nil {
			JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		insertedTodo, err := tds.AddTodo(todo.Name)
		if err != nil {
			JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(insertedTodo)

	})
	r.Put("/api/v1/todo/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		var todo UpdatedTodo
		err = json.NewDecoder(r.Body).Decode(&todo)
		if err != nil {
			JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		res, err := db.Exec("UPDATE todos SET name = ?, completed = ? WHERE id = ?", todo.Name, todo.Completed, id)
		if err != nil {
			JSONError(w, err.Error(), http.StatusInternalServerError)
		}
		rowsUpdated, err := res.RowsAffected()
		if rowsUpdated == 1 {
			json.NewEncoder(w).Encode(Todo{
				ID:        idInt,
				Name:      todo.Name,
				Completed: todo.Completed,
			})
			return
		}
		JSONError(w, "something went wrong", http.StatusBadRequest)
	})
	r.Delete("/todo/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		_, err := tds.RemoveTodo(id)
		if err != nil {
			errDiv().Render(context.Background(), w)
			return
		}
	})
	r.Delete("/api/v1/todo/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		rowsDeleted, err := tds.RemoveTodo(id)
		if err != nil {
			JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]int64{"deletedCount": rowsDeleted})

	})
	r.Delete("/api/v1/todo", func(w http.ResponseWriter, r *http.Request) {
		response := "delete all todo"
		w.Write([]byte(response))
	})
	r.Get("/api/v1/todo/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		todo, err := tds.GetTodo(id)
		if err != nil {
			JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(todo)
	})
	r.Get("/api/v1/todo", func(w http.ResponseWriter, r *http.Request) {
		todos, err := tds.GetAllTodos()
		if err != nil {
			JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(todos)
	})

	http.ListenAndServe(":8080", r)
}
