package storage

import (
	"log"
	"log/slog"
	"os"
	"reflect"
	"testing"
	"todoapp/2/storage"
)

var logger *slog.Logger

func TestMain(m *testing.M) {

	// Setup: This runs before any test
	log.Println("Setting up test environment")

	// Run the tests
	exitCode := m.Run()

	// Teardown: This runs after all tests
	log.Println("Cleaning up test environment")

	// Exit with the test status
	os.Exit(exitCode)
}

// TestAddTodo tests the AddTodo and TestGetTodoByID functions
func TestAddTodo(t *testing.T) {
	tl := storage.NewTodoList()
	addedTodo := tl.AddTodo("Test Todo")
	// Check if the todo was added successfully
	storedTodo, exists := tl.GetTodoByID(addedTodo.ID)
	if !exists {
		t.Fatalf("Expected todo to be added but was not found")
	}

	if storedTodo.Description != "Test Todo" || storedTodo.Completed != false {
		t.Errorf("Added todo has incorrect values. Got %+v, expected %+v", storedTodo, addedTodo)
	}
}

// TestGetAllTodos tests if GetAllTodos returns all added todos
func TestGetAllTodos(t *testing.T) {
	// Add multiple todos
	tl := storage.NewTodoList()
	addedTodo1 := tl.AddTodo("Test Todo 1")
	addedTodo2 := tl.AddTodo("Test Todo 2")

	todoList := tl.GetAllTodos()

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
	tl := storage.NewTodoList()
	todo := tl.AddTodo("Todo to Update")

	// Update the todo
	newTodoForUpdate := storage.Todo{Description: "Updated Todo", Completed: true, ID: todo.ID}
	updatedTodo, success := tl.UpdateTodoByID(todo.ID, newTodoForUpdate)

	if !success {
		t.Fatalf("Failed to update todo with ID %d", todo.ID)
	}

	retrievedTodo, _ := tl.GetTodoByID(todo.ID)
	if !reflect.DeepEqual(retrievedTodo, updatedTodo) {
		t.Errorf("Updated todo does not match. Got %+v, expected %+v", retrievedTodo, updatedTodo)
	}
}

// TestDeleteTodoByID tests the DeleteTodoByID function
func TestDeleteTodoByID(t *testing.T) {
	tl := storage.NewTodoList()
	addedTodo := tl.AddTodo("Todo to Delete")
	success := tl.DeleteTodoByID(addedTodo.ID)
	if !success {
		t.Fatalf("Failed to delete todo with ID %d", addedTodo.ID)
	}

	_, exists := tl.GetTodoByID(addedTodo.ID)
	if exists {
		t.Errorf("Expected todo with ID %d to be deleted, but it still exists", addedTodo.ID)
	}
}
