package unit_test

import (
	"context"
	"log"
	"os"
	"testing"
	"todoapp/3/storage"

	"github.com/stretchr/testify/assert"
)

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
	err := todoList.UpdateTodoByID(ctx, todo.ID, newTodoForUpdate)
	if err != nil {
		t.Fatalf("Failed to update todo with ID %d", todo.ID)
	}
	retrievedTodo, _ := todoList.GetTodoByID(ctx, todo.ID)
	assert.Equal(t, newTodoForUpdate, retrievedTodo)
}

// TestDeleteTodoByID tests the DeleteTodoByID function
func TestDeleteTodoByID(t *testing.T) {
	store := storage.NewInMemoryStore()
	todoList := storage.NewTodoListWithOptions(storage.Options{
		Store: store,
	})
	ctx := context.Background()
	todo, _ := todoList.AddTodo(ctx, "Todo to Delete")

	err := todoList.DeleteTodoByID(ctx, todo.ID)
	if err != nil {
		t.Fatalf("Failed to delete todo with ID %d", todo.ID)
	}

	todo, _ = todoList.GetTodoByID(ctx, todo.ID)
	assert.Nil(t, todo)
}
