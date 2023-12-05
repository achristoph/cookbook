package storage

import (
	"io"
	"log/slog"
	"os"
	"sync"
)

// Todo struct represents a task with an ID and a description
type Todo struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type TodoList struct {
	Logger    *slog.Logger
	todos     sync.Map   // In-memory storage for todos using sync.Map
	mu        sync.Mutex // Ensures safe ID increments
	idCounter int        // ID counter
}

type Options struct {
	Logger *slog.Logger
}

// TodoList represents a set of todos
// It includes relevant external service such as logger instance
// the todo functions are now converted to methods, this allows using logging without having pass the logger to each method
// NewTodoList creates a TodoList with a default logger
func NewTodoList() *TodoList {
	return NewTodoListWithOptions(Options{})
}

// NewTodoListWithOptions allows custom logger configuration
func NewTodoListWithOptions(options Options) *TodoList {
	logger := options.Logger
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil)) // Default logger
	}
	return &TodoList{
		Logger:    logger,
		idCounter: 1,
	}
}

// Using constructor to create a todo
func NewTodo(description string) *Todo {
	return &Todo{Description: description}
}

// GetAllTodos retrieves all todos from the storage
func (t *TodoList) GetAllTodos() []*Todo {
	var todoList []*Todo

	t.todos.Range(func(key, value interface{}) bool {
		todoList = append(todoList, value.(*Todo))
		return true
	})
	t.Logger.Info("Listing all todos")
	return todoList
}

// AddTodo adds a new todo to the storage
// Return the inserted todo with the generated ID
func (t *TodoList) AddTodo(description string) *Todo {
	t.mu.Lock()
	defer t.mu.Unlock()
	newTodo := &Todo{ID: t.idCounter, Description: description}
	t.idCounter++
	t.todos.Store(newTodo.ID, newTodo)
	t.Logger.Info("Added a todo")
	return newTodo
}

// GetTodoByID retrieves a todo by its ID
// Return false if it does not exist or fails to cast
func (t *TodoList) GetTodoByID(id int) (*Todo, bool) {
	obj, exists := t.todos.Load(id)
	if !exists {
		return nil, false
	}
	todo, ok := obj.(*Todo)
	if !ok {
		return nil, false
	}
	return todo, true
}

// UpdateTodoByID updates a todo by its ID
// The existing todo is replaced with a new one using the same id
// Return false if the todo with the given ID does not exist
func (t *TodoList) UpdateTodoByID(id int, updatedTodo Todo) (*Todo, bool) {
	if _, exists := t.todos.Load(id); !exists {
		return nil, false
	}
	t.todos.Store(id, updatedTodo)
	return nil, true
}

// DeleteTodoByID deletes a todo by its ID
// Return false if the todo with the given ID does not exist
func (t *TodoList) DeleteTodoByID(id int) bool {
	if _, exists := t.todos.Load(id); !exists {
		return false
	}
	t.todos.Delete(id)
	return true
}

// Disable logging by setting output to io.Discard
func (t *TodoList) DisableLogging() {
	t.Logger = slog.New(slog.NewTextHandler(io.Discard, nil))
}
