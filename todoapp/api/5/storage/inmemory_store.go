package storage

import (
	"context"
	"sync"
)

// InMemoryStore is a thread-safe in-memory implementation of TodoStore.
type InMemoryStore struct {
	todos     sync.Map   // Stores todos using their ID as the key
	mu        sync.Mutex // Ensures safe ID increments
	idCounter int        // ID counter for generating unique IDs
}

// NewInMemoryStore creates an in-memory storage instance.
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{idCounter: 1}
}

// AddTodo adds a new todo to the in-memory store.
func (s *InMemoryStore) AddTodo(ctx context.Context, description string) (*Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check for duplicate description
	s.todos.Range(func(_, value interface{}) bool {
		if todo := value.(*Todo); todo.Description == description {
			return false
		}
		return true
	})

	// Assign unique ID
	todo := &Todo{ID: int(s.idCounter), Description: description, Completed: false}
	s.todos.Store(s.idCounter, todo)
	s.idCounter++
	return todo, nil
}

// GetAllTodos retrieves all todos.
func (s *InMemoryStore) GetAllTodos(ctx context.Context) ([]*Todo, error) {
	var todoList []*Todo
	s.todos.Range(func(_, value interface{}) bool {
		todoList = append(todoList, value.(*Todo))
		return true
	})
	return todoList, nil
}

// GetTodoByID retrieves a todo by ID.
func (s *InMemoryStore) GetTodoByID(ctx context.Context, id int) (*Todo, error) {
	if value, exists := s.todos.Load(id); exists {
		return value.(*Todo), nil
	}
	return nil, NewTodoNotFoundError(id)
}

// UpdateTodoByID updates an existing todo.
func (s *InMemoryStore) UpdateTodoByID(ctx context.Context, id int, updatedTodo *Todo) error {
	if _, exists := s.todos.Load(id); !exists {
		return NewTodoNotFoundError(id)
	}
	s.todos.Store(id, updatedTodo)
	return nil
}

// DeleteTodoByID removes a todo by ID.
func (s *InMemoryStore) DeleteTodoByID(ctx context.Context, id int) error {
	if _, exists := s.todos.Load(id); !exists {
		return NewTodoNotFoundError(id)
	}
	s.todos.Delete(id)
	return nil
}
