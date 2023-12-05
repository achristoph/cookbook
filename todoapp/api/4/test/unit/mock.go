package unit_test

import (
	"io"
	"os"

	"github.com/stretchr/testify/mock"
)

// MockStorageIO defines a mock version of StorageIO.
type MockStorageIO struct {
	mock.Mock
}

func (m *MockStorageIO) OpenFile(path string, flag int, perm os.FileMode) (*os.File, error) {
	args := m.Called(path, flag, perm)
	return args.Get(0).(*os.File), args.Error(1)
}

func (m *MockStorageIO) CreateFile(path string) (*os.File, error) {
	args := m.Called(path)
	return args.Get(0).(*os.File), args.Error(1)
}

func (m *MockStorageIO) ReadFile(path string) ([]byte, error) {
	args := m.Called(path)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockStorageIO) WriteFile(path string, data []byte) error {
	args := m.Called(path, data)
	return args.Error(0)
}

func (m *MockStorageIO) EncodeJSON(writer io.Writer, data interface{}) error {
	args := m.Called(writer, data)
	return args.Error(0)
}

func (m *MockStorageIO) DecodeJSON(reader io.Reader, out interface{}) error {
	args := m.Called(reader, out)
	return args.Error(0)
}
