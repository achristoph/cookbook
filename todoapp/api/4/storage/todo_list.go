package storage

import (
	"context"
	"io"
	"log/slog"
	"os"
)

// TodoList represents a set of todos backed by a storage implementation.
type TodoList struct {
	Logger    *slog.Logger
	Store     TodoStore // Can be SQLite or InMemoryStore
	StorageIO StorageIOInterface
}

// Todo struct represents a task with an ID and a description
type Todo struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type Options struct {
	Logger *slog.Logger
	Store  TodoStore
}

// TodoList represents a set of todos
// It includes relevant external service such as logger instance
// the todo functions are now converted to methods, this allows using logging without having pass the logger to each method
// NewTodoList creates a TodoList with a default logger
func NewTodoList() *TodoList {
	return NewTodoListWithOptions(Options{})
}

// Using constructor to create a todo
func NewTodo(description string) *Todo {
	return &Todo{Description: description}
}

// NewTodoListWithOptions creates a TodoList with custom storage.
func NewTodoListWithOptions(options Options) *TodoList {
	if options.Store == nil {
		options.Store = NewInMemoryStore() // Assume you have an in-memory implementation
	}
	if options.Logger == nil {
		options.Logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	return &TodoList{
		Logger:    options.Logger,
		Store:     options.Store,
		StorageIO: &StorageIO{}, // Default StorageIO instance
	}
}

// AddTodo adds a new todo using the configured storage backend.
func (t *TodoList) AddTodo(ctx context.Context, description string) (*Todo, error) {
	todo, err := t.Store.AddTodo(ctx, description)
	if err != nil {
		t.Logger.Error("Failed to add todo", "error", err)
		return nil, err
	}
	t.Logger.Info("Added a todo", "id", todo.ID)
	return todo, nil

}

// GetAllTodos retrieves all todos.
func (t *TodoList) GetAllTodos(ctx context.Context) ([]*Todo, error) {
	t.Logger.Info("Listing all todos")
	todos, err := t.Store.GetAllTodos(ctx)
	return todos, err
}

// GetTodoByID retrieves a todo by ID.
func (t *TodoList) GetTodoByID(ctx context.Context, id int) (*Todo, error) {
	t.Logger.Info("Getting a todo", "id", id)
	todo, err := t.Store.GetTodoByID(ctx, id)
	return todo, err
}

// UpdateTodoByID updates a todo by its ID.
func (t *TodoList) UpdateTodoByID(ctx context.Context, id int, updatedTodo *Todo) error {
	t.Logger.Info("Updating a todo", "id", id)
	return t.Store.UpdateTodoByID(ctx, id, updatedTodo)
}

// DeleteTodoByID deletes a todo by ID.
func (t *TodoList) DeleteTodoByID(ctx context.Context, id int) error {
	t.Logger.Info("Deleting a todo", "id", id)
	return t.Store.DeleteTodoByID(ctx, id)
}

func (t *TodoList) Download(ctx context.Context, path string) error {
	t.Logger.Info("Downloading todos to file", "path", path)

	todos, err := t.Store.GetAllTodos(ctx)
	if err != nil {
		t.Logger.Error("Failed to retrieve todos for download", "error", err)
		return err
	}

	file, err := t.StorageIO.CreateFile(path)
	if err != nil {
		t.Logger.Error("Failed to create file", "path", path, "error", err)
		return err
	}
	defer file.Close()

	if err := t.StorageIO.EncodeJSON(file, todos); err != nil {
		t.Logger.Error("Failed to write todos to file", "error", err)
		return err
	}

	t.Logger.Info("Successfully downloaded todos", "path", path)
	return nil
}

func (t *TodoList) Upload(ctx context.Context, path string) error {
	t.Logger.Info("Uploading todos from file", "path", path)

	file, err := t.StorageIO.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		t.Logger.Error("Failed to open file", "path", path, "error", err)
		return err
	}
	defer file.Close()

	var todos []*Todo
	if err := t.StorageIO.DecodeJSON(file, &todos); err != nil {
		t.Logger.Error("Failed to parse todos from file", "error", err)
		return err
	}

	for _, todo := range todos {
		_, err := t.Store.AddTodo(ctx, todo.Description)
		if err != nil {
			t.Logger.Error("Failed to add todo", "id", todo.ID, "error", err)
		}
	}

	t.Logger.Info("Successfully uploaded todos from file", "path", path)
	return nil
}

// Disable logging by setting output to io.Discard
func (t *TodoList) DisableLogging() {
	t.Logger = slog.New(slog.NewTextHandler(io.Discard, nil))
}
