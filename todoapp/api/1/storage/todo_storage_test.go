package storage

import (
	"reflect"
	"sync"
	"testing"
)

// TestAddTodo tests the AddTodo and TestGetTodoByID functions
func TestAddTodo(t *testing.T) {
	todo := &Todo{Description: "Test Todo", Completed: false}

	addedTodo := AddTodo(todo)

	// Check if the todo was added successfully
	storedTodo, exists := GetTodoByID(addedTodo.ID)
	if !exists {
		t.Fatalf("Expected todo to be added but was not found")
	}

	if storedTodo.Description != "Test Todo" || storedTodo.Completed != false {
		t.Errorf("Added todo has incorrect values. Got %+v, expected %+v", storedTodo, todo)
	}
}

// TestGetAllTodos tests if GetAllTodos returns all added todos
func TestGetAllTodos(t *testing.T) {
	// Reset the storage for the test
	todos = sync.Map{}
	idCounter = 1

	// Add multiple todos
	todo1 := &Todo{Description: "Test Todo 1", Completed: false}
	todo2 := &Todo{Description: "Test Todo 2", Completed: true}
	addedTodo1 := AddTodo(todo1)
	addedTodo2 := AddTodo(todo2)

	todoList := GetAllTodos()

	if len(todoList) != 2 {
		t.Fatalf("Expected 2 todos, but got %d", len(todoList))
	}

	// Check if todos returned match what was added
	if !reflect.DeepEqual(todoList[0], addedTodo1) && !reflect.DeepEqual(todoList[1], addedTodo2) {
		t.Errorf("Todos returned do not match what was added. Got %+v", todoList)
	}
}

// TestUpdateTodoByID tests if the UpdateTodoByID function works correctly
func TestUpdateTodoByID(t *testing.T) {
	todo := &Todo{Description: "Todo to Update", Completed: false}

	AddTodo(todo)

	// Update the todo
	newTodoForUpdate := Todo{Description: "Updated Todo", Completed: true, ID: todo.ID}
	updatedTodo, success := UpdateTodoByID(todo.ID, newTodoForUpdate)

	if !success {
		t.Fatalf("Failed to update todo with ID %d", todo.ID)
	}

	retrievedTodo, _ := GetTodoByID(todo.ID)
	if !reflect.DeepEqual(retrievedTodo, updatedTodo) {
		t.Errorf("Updated todo does not match. Got %+v, expected %+v", retrievedTodo, updatedTodo)
	}
}

// TestDeleteTodoByID tests the DeleteTodoByID function
func TestDeleteTodoByID(t *testing.T) {
	todo := &Todo{Description: "Todo to Delete", Completed: false}

	addedTodo := AddTodo(todo)

	success := DeleteTodoByID(addedTodo.ID)
	if !success {
		t.Fatalf("Failed to delete todo with ID %d", addedTodo.ID)
	}

	_, exists := GetTodoByID(addedTodo.ID)
	if exists {
		t.Errorf("Expected todo with ID %d to be deleted, but it still exists", addedTodo.ID)
	}
}
