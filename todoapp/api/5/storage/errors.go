package storage

import "fmt"

// Error types for todo operations
type TodoError struct {
	Code    string
	Message string
	Err     error
}

func (e *TodoError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.Err.Error())
	}
	return e.Message
}

// Error codes
const (
	ErrTodoNotFound     = "TODO_NOT_FOUND"
	ErrInvalidInput     = "INVALID_INPUT"
	ErrStorageError     = "STORAGE_ERROR"
	ErrDuplicateTodo    = "DUPLICATE_TODO"
	ErrOperationTimeout = "OPERATION_TIMEOUT"
)

// Helper functions to create specific errors
func NewTodoNotFoundError(id int) *TodoError {
	return &TodoError{
		Code:    ErrTodoNotFound,
		Message: fmt.Sprintf("Todo with ID %d not found", id),
	}
}

func NewInvalidInputError(message string) *TodoError {
	return &TodoError{
		Code:    ErrInvalidInput,
		Message: message,
	}
}

func NewStorageError(err error) *TodoError {
	return &TodoError{
		Code:    ErrStorageError,
		Message: "Storage operation failed",
		Err:     err,
	}
}

func NewDuplicateTodoError(description string) *TodoError {
	return &TodoError{
		Code:    ErrDuplicateTodo,
		Message: fmt.Sprintf("Todo with description '%s' already exists", description),
	}
}

func NewOperationTimeoutError() *TodoError {
	return &TodoError{
		Code:    ErrOperationTimeout,
		Message: "Operation timed out",
	}
}
