package unit_test

import (
	"context"
	"log"
	"log/slog"
	"os"
	"testing"
	"todoapp/3/storage"

	"github.com/stretchr/testify/assert"
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
	store := storage.NewInMemoryStore()
	todoList := storage.NewTodoListWithOptions(storage.Options{
		Store: store,
	})
	ctx := context.Background()
	// Add a new todo
	todo, err := todoList.AddTodo(ctx, "Test in-memory todo")
	assert.NoError(t, err)
	assert.NotNil(t, todo)
	assert.Equal(t, "Test in-memory todo", todo.Description)

	// Retrieve and verify
	fetchedTodo, err := todoList.GetTodoByID(ctx, todo.ID)
	assert.NoError(t, err)
	assert.Equal(t, todo.ID, fetchedTodo.ID)
}

// TestGetAllTodos tests if GetAllTodos returns all added todos
func TestGetAllTodos(t *testing.T) {
	store := storage.NewInMemoryStore()
	todoList := storage.NewTodoListWithOptions(storage.Options{
		Store: store,
	})
	ctx := context.Background()
	todoList.AddTodo(ctx, "Test Todo 1")
	todoList.AddTodo(ctx, "Test Todo 2")
	allTodos, _ := todoList.GetAllTodos(ctx)

	assert.Equal(t, 2, len(allTodos))
}

// TestUpdateTodoByID tests if the UpdateTodoByID function works correctly
func TestUpdateTodoByID(t *testing.T) {
	store := storage.NewInMemoryStore()
	todoList := storage.NewTodoListWithOptions(storage.Options{
		Store: store,
	})
	ctx := context.Background()

	todo, _ := todoList.AddTodo(ctx, "Todo to Update")

	// Update the todo
	newTodoForUpdate := &storage.Todo{ID: todo.ID, Description: "Updated Todo", Completed: true}
	todoList.UpdateTodoByID(ctx, todo.ID, newTodoForUpdate)

	retrievedTodo, _ := todoList.GetTodoByID(ctx, todo.ID)
	assert.Equal(t, newTodoForUpdate, retrievedTodo)
	// if !reflect.DeepEqual(retrievedTodo, updatedTodo.ID) {
	// 	t.Errorf("Updated todo does not match. Got %+v, expected %+v", retrievedTodo, updatedTodo)
	// }
}

// // TestDeleteTodoByID tests the DeleteTodoByID function
// func TestDeleteTodoByID(t *testing.T) {
// 	tl := NewTodoList()
// 	addedTodo := tl.AddTodo("Todo to Delete")
// 	success := tl.DeleteTodoByID(addedTodo.ID)
// 	if !success {
// 		t.Fatalf("Failed to delete todo with ID %d", addedTodo.ID)
// 	}

// 	_, exists := tl.GetTodoByID(addedTodo.ID)
// 	if exists {
// 		t.Errorf("Expected todo with ID %d to be deleted, but it still exists", addedTodo.ID)
// 	}
// }
