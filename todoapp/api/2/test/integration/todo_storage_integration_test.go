package storage

import (
	"log"
	"log/slog"
	"os"
	"testing"
	"todoapp/2/storage"
)

func TestIntegration_TodoListWithLogFile(t *testing.T) {
	file, err := os.OpenFile("todo.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	logger := slog.New(slog.NewTextHandler(file, nil))
	todoList := storage.NewTodoListWithOptions(storage.Options{Logger: logger})
	todoList.AddTodo("item 1")
	todos := todoList.GetAllTodos()
	log.Println(len(todos))
	file.Close()
}
