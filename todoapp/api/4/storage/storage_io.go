package storage

import (
	"encoding/json"
	"io"
	"os"
)

type StorageIOInterface interface {
	OpenFile(path string, flag int, perm os.FileMode) (*os.File, error)
	CreateFile(path string) (*os.File, error)
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte) error
	EncodeJSON(writer io.Writer, data interface{}) error
	DecodeJSON(reader io.Reader, out interface{}) error
}

// StorageIO abstracts file I/O operations for easier testing.
type StorageIO struct{}

// OpenFile opens a file with the given path and flag.
func (s *StorageIO) OpenFile(path string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(path, flag, perm)
}

// CreateFile creates a new file at the given path.
func (s *StorageIO) CreateFile(path string) (*os.File, error) {
	return os.Create(path)
}

// ReadFile reads a file and returns its contents.
func (s *StorageIO) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// WriteFile writes data to a file.
func (s *StorageIO) WriteFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

// EncodeJSON encodes data into JSON format and writes to the given file.
func (s *StorageIO) EncodeJSON(writer io.Writer, data interface{}) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ") // Pretty-print JSON
	return encoder.Encode(data)
}

// DecodeJSON decodes JSON data from the given file.
func (s *StorageIO) DecodeJSON(reader io.Reader, v interface{}) error {
	return json.NewDecoder(reader).Decode(v)
}
