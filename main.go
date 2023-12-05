package main

import (
	"log/slog"
	"os"
)

func main() {
	// Create a logger that writes to stdout
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Log messages with different levels
	logger.Info("This is an info message")
	logger.Warn("This is a warning message")
	logger.Error("This is an error message")
}
