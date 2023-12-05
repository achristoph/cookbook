package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"todoapp/2/storage"
)

var todoList *storage.TodoList

func getTodosHandler(w http.ResponseWriter, _ *http.Request) {
	list := todoList.GetAllTodos()

	// Set response header
	w.Header().Set("Content-Type", "application/json")

	// Attempt to encode the list into JSON
	if err := json.NewEncoder(w).Encode(list); err != nil {
		// If an error occurs during encoding, return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to encode todos",
		})
		return
	}
}

func createTodoHandler(w http.ResponseWriter, r *http.Request) {
	var newTodo *storage.Todo
	if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	todoList.AddTodo(newTodo.Description)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTodo)
}

func getTodoHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL
	idStr := r.URL.Path[len("/todos/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Retrieve the todo by ID
	todo, exists := todoList.GetTodoByID(id)
	if !exists {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Todo not found",
		})
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
	todo, exists := todoList.UpdateTodoByID(id, updatedTodo)
	if !exists {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Failed to update a todo with id %d", id),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
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
	exists := todoList.DeleteTodoByID(id)
	if !exists {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Failed to delete a todo with id %d", id),
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	todoList = storage.NewTodoList()
	// Create a new ServeMux to handle routes
	mux := http.NewServeMux()

	// Route for getting all todos and creating a new todo
	mux.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getTodosHandler(w, r)
		case http.MethodPost:
			createTodoHandler(w, r)
		default:
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
