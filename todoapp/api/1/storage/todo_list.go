package storage

import (
	"sync"
)

// Todo struct represents a task with an ID and a description
type Todo struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

// In-memory storage for todos using sync.Map
var todos sync.Map
var idCounter = 1
var counterLock sync.Mutex // To safely increment the ID counter

// GetAllTodos retrieves all todos from the storage
func GetAllTodos() []*Todo {
	var todoList []*Todo

	todos.Range(func(key, value interface{}) bool {
		todoList = append(todoList, value.(*Todo))
		return true
	})

	return todoList
}

// AddTodo adds a new todo to the storage
// Return the inserted todo with the generated ID
func AddTodo(newTodo *Todo) *Todo {
	counterLock.Lock()
	newTodo.ID = idCounter
	idCounter++
	counterLock.Unlock()
	todos.Store(newTodo.ID, newTodo)
	return newTodo
}

// GetTodoByID retrieves a todo by its ID
// Return false if it does not exist or fails to cast
func GetTodoByID(id int) (*Todo, bool) {
	obj, exists := todos.Load(id)
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
func UpdateTodoByID(id int, updatedTodo Todo) (*Todo, bool) {
	if _, exists := todos.Load(id); !exists {
		return nil, false
	}
	todos.Store(id, updatedTodo)
	return nil, true
}

// DeleteTodoByID deletes a todo by its ID
// Return false if the todo with the given ID does not exist
func DeleteTodoByID(id int) bool {
	if _, exists := todos.Load(id); !exists {
		return false
	}
	todos.Delete(id)
	return true
}
