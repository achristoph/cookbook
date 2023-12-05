package storage

import "context"

// TodoStore defines storage operations for todos.
type TodoStore interface {
	AddTodo(ctx context.Context, description string) (*Todo, error)
	GetAllTodos(ctx context.Context) ([]*Todo, error)
	GetTodoByID(ctx context.Context, id int) (*Todo, error)
	UpdateTodoByID(ctx context.Context, id int, updatedTodo *Todo) error
	DeleteTodoByID(ctx context.Context, id int) error
}
