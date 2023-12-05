package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"todoapp/1/storage"
)

func getTodosHandler(w http.ResponseWriter, r *http.Request) {
	// Get all todos from the storage package
	todoList := storage.GetAllTodos()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todoList)
}

func createTodoHandler(w http.ResponseWriter, r *http.Request) {
	var newTodo *storage.Todo
	if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Add new todo to storage
	storage.AddTodo(newTodo)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTodo)
}

func getTodoHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL
	idStr := r.URL.Path[len("/todos/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Todo ID", http.StatusBadRequest)
		return
	}

	// Retrieve the todo by ID
	todo, exists := storage.GetTodoByID(id)
	if !exists {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func updateTodoHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL
	idStr := r.URL.Path[len("/todos/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Todo ID", http.StatusBadRequest)
		return
	}

	var updatedTodo storage.Todo
	if err := json.NewDecoder(r.Body).Decode(&updatedTodo); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Update the todo in storage
	updatedTodo.ID = id
	_, exist := storage.UpdateTodoByID(id, updatedTodo)
	if !exist {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTodo)
}

func deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL
	idStr := r.URL.Path[len("/todos/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Todo ID", http.StatusBadRequest)
		return
	}

	// Delete the todo by ID
	if !storage.DeleteTodoByID(id) {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	// Create a new ServeMux to handle routes
	mux := http.NewServeMux()

	// Route for getting all todos and creating a new todo
	mux.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getTodosHandler(w, r)
		} else if r.Method == http.MethodPost {
			createTodoHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Route for handling specific todo by ID (get, update, delete)
	mux.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getTodoHandler(w, r)
		case http.MethodPut:
			updateTodoHandler(w, r)
		case http.MethodDelete:
			deleteTodoHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Start the HTTP server
	fmt.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
