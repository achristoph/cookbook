package storage

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteTodoStore implements TodoStore using SQLite
type SQLiteTodoStore struct {
	DB *sql.DB
}

// NewSQLiteTodoStore initializes a SQLite-backed store
func NewSQLiteTodoStore(db *sql.DB) *SQLiteTodoStore {
	return &SQLiteTodoStore{DB: db}
}

// AddTodo inserts a new todo
func (s *SQLiteTodoStore) AddTodo(ctx context.Context, description string) (*Todo, error) {
	query := "INSERT INTO todos (description, completed) VALUES (?, ?)"
	result, err := s.DB.ExecContext(ctx, query, description, false)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Todo{ID: int(id), Description: description, Completed: false}, nil
}

// GetAllTodos fetches all todos
func (s *SQLiteTodoStore) GetAllTodos(ctx context.Context) ([]*Todo, error) {
	rows, err := s.DB.QueryContext(ctx, "SELECT id, description, completed FROM todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*Todo
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.ID, &todo.Description, &todo.Completed); err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}
	return todos, nil
}

// GetTodoByID fetches a todo by ID
func (s *SQLiteTodoStore) GetTodoByID(ctx context.Context, id int) (*Todo, error) {
	query := "SELECT id, description, completed FROM todos WHERE id = ?"
	row := s.DB.QueryRowContext(ctx, query, id)

	var todo Todo
	if err := row.Scan(&todo.ID, &todo.Description, &todo.Completed); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &todo, nil
}

// UpdateTodoByID updates a todo
func (s *SQLiteTodoStore) UpdateTodoByID(ctx context.Context, id int, updatedTodo *Todo) error {
	query := "UPDATE todos SET description = ?, completed = ? WHERE id = ?"
	_, err := s.DB.ExecContext(ctx, query, updatedTodo.Description, updatedTodo.Completed, id)
	return err
}

// DeleteTodoByID deletes a todo
func (s *SQLiteTodoStore) DeleteTodoByID(ctx context.Context, id int) error {
	query := "DELETE FROM todos WHERE id = ?"
	_, err := s.DB.ExecContext(ctx, query, id)
	return err
}
