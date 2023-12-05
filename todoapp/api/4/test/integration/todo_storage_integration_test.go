package integration_test

import (
	"context"
	"database/sql"
	"testing"
	"todoapp/3/storage"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestTodoList_AddTodo(t *testing.T) {
	// Use in-memory SQLite
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	defer db.Close()

	// Create table
	_, err = db.Exec(`CREATE TABLE todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		description TEXT NOT NULL,
		completed BOOLEAN NOT NULL
	)`)
	assert.NoError(t, err)

	// Initialize SQLite storage
	store := storage.NewSQLiteTodoStore(db)

	// Create TodoList with SQLite backend
	todoList := storage.NewTodoListWithOptions(storage.Options{
		Store: store,
	})

	// Add a new todo
	ctx := context.Background()
	todo, err := todoList.AddTodo(ctx, "Write integration test")
	assert.NoError(t, err)
	assert.NotNil(t, todo)
	assert.Equal(t, "Write integration test", todo.Description)

	// Fetch the todo and verify
	fetchedTodo, err := todoList.GetTodoByID(ctx, todo.ID)
	assert.NoError(t, err)
	assert.NotNil(t, fetchedTodo)
	assert.Equal(t, todo.ID, fetchedTodo.ID)
	assert.Equal(t, "Write integration test", fetchedTodo.Description)
}
