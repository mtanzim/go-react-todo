package main

import (
	"database/sql"
)

type todoService struct {
	db *sql.DB
}

func NewTodoService(db *sql.DB) *todoService {
	return &todoService{db}
}

func (tdl *todoService) GetTodo(id string) (Todo, error) {
	row := tdl.db.QueryRow("SELECT id, name, completed FROM todos WHERE id = ?", id)
	var todo Todo
	err := row.Scan(&todo.ID, &todo.Name, &todo.Completed)
	if err != nil {
		return Todo{}, err
	}
	return todo, nil
}

func (tdl *todoService) AddTodo(title string) (Todo, error) {

	res, err := tdl.db.Exec("INSERT INTO todos (name, completed) VALUES (?, ?)", title, false)
	if err != nil {
		return Todo{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return Todo{}, err
	}
	insertedTodo := Todo{
		ID:        id,
		Name:      title,
		Completed: false,
	}
	return insertedTodo, nil
}

func (tdl *todoService) RemoveTodo(id string) (int64, error) {

	res, err := tdl.db.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		return -1, err
	}
	rowsDeleted, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}
	return rowsDeleted, nil
}

func (tdl *todoService) GetAllTodos() ([]Todo, error) {
	rows, err := tdl.db.Query("SELECT id, name, completed FROM todos ORDER BY id DESC")
	if err != nil {
		return []Todo{}, err
	}
	todos := []Todo{}
	defer rows.Close()
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Name, &todo.Completed)
		if err != nil {
			return []Todo{}, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}
