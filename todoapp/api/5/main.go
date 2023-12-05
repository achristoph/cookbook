package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
	"todoapp/5/storage"
)

var todoList *storage.TodoList

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func writeErrorResponse(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	var response ErrorResponse
	if todoErr, ok := err.(*storage.TodoError); ok {
		response = ErrorResponse{
			Error:   todoErr.Error(),
			Code:    todoErr.Code,
			Message: todoErr.Message,
		}
	} else {
		response = ErrorResponse{
			Error: err.Error(),
		}
	}

	json.NewEncoder(w).Encode(response)
}

func getTodosHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	list, err := todoList.GetAllTodos(ctx)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, storage.NewStorageError(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func createTodoHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var newTodo *storage.Todo
	if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, storage.NewInvalidInputError("Invalid request payload"))
		return
	}

	if newTodo.Description == "" {
		writeErrorResponse(w, http.StatusBadRequest, storage.NewInvalidInputError("Todo description cannot be empty"))
		return
	}

	todo, err := todoList.AddTodo(ctx, newTodo.Description)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, storage.NewStorageError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

func getTodoHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/todos/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, storage.NewInvalidInputError("Invalid todo ID"))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	todo, err := todoList.GetTodoByID(ctx, id)
	if err != nil {
		if err.Error() == "todo not found" {
			writeErrorResponse(w, http.StatusNotFound, storage.NewTodoNotFoundError(id))
		} else {
			writeErrorResponse(w, http.StatusInternalServerError, storage.NewStorageError(err))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func updateTodoHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/todos/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, storage.NewInvalidInputError("Invalid todo ID"))
		return
	}

	var updatedTodo *storage.Todo
	if err := json.NewDecoder(r.Body).Decode(&updatedTodo); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, storage.NewInvalidInputError("Invalid request payload"))
		return
	}

	if updatedTodo.Description == "" {
		writeErrorResponse(w, http.StatusBadRequest, storage.NewInvalidInputError("Todo description cannot be empty"))
		return
	}

	updatedTodo.ID = id
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	err = todoList.UpdateTodoByID(ctx, id, updatedTodo)
	if err != nil {
		if err.Error() == "todo not found" {
			writeErrorResponse(w, http.StatusNotFound, storage.NewTodoNotFoundError(id))
		} else {
			writeErrorResponse(w, http.StatusInternalServerError, storage.NewStorageError(err))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTodo)
}

func deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/todos/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, storage.NewInvalidInputError("Invalid todo ID"))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	err = todoList.DeleteTodoByID(ctx, id)
	if err != nil {
		if err.Error() == "todo not found" {
			writeErrorResponse(w, http.StatusNotFound, storage.NewTodoNotFoundError(id))
		} else {
			writeErrorResponse(w, http.StatusInternalServerError, storage.NewStorageError(err))
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func downloadTodosHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	filename := "todos.json"
	err := todoList.Download(ctx, filename)
	if err != nil {
		http.Error(w, "Failed to download todos", http.StatusInternalServerError)
		return
	}

	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, "Error opening file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "application/json")
	http.ServeFile(w, r, filename)
}

func uploadTodosHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var file io.Reader
	var err error

	// Check for multipart file upload
	if r.Header.Get("Content-Type") == "multipart/form-data" {
		err = r.ParseMultipartForm(10 << 20) // 10MB limit
		if err != nil {
			http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
			return
		}

		uploadedFile, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Failed to retrieve uploaded file", http.StatusBadRequest)
			return
		}
		defer uploadedFile.Close()

		file = uploadedFile
	} else {
		// Fallback to JSON body with "path"
		var requestData struct {
			Path string `json:"path"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		file, err = os.Open(requestData.Path)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to open file: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		defer file.(*os.File).Close()
	}

	// Decode todos
	var todos []*storage.Todo
	if err := todoList.StorageIO.DecodeJSON(file, &todos); err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse todos: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Add todos to storage
	for _, todo := range todos {
		_, err := todoList.Store.AddTodo(ctx, todo.Description)
		if err != nil {
			todoList.Logger.Error("Failed to add todo", "id", todo.ID, "error", err)
		}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Todos uploaded successfully",
	})
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

	mux.HandleFunc("/todos/download", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		downloadTodosHandler(w, r)
	})

	mux.HandleFunc("/todos/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		uploadTodosHandler(w, r)
	})

	// Start the HTTP server
	fmt.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
