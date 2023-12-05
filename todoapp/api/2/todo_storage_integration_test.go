package storage

import (
	"log"
	"log/slog"
	"os"
	"testing"
)

func TestIntegration_TodoListWithLogFile(t *testing.T) {
	file, err := os.OpenFile("todo.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	logger = slog.New(slog.NewTextHandler(file, nil))
	todoList := NewTodoListWithOptions(Options{Logger: logger})
	todoList.AddTodo("item 1")
	todos := todoList.GetAllTodos()
	log.Println(len(todos))
	file.Close()
}
