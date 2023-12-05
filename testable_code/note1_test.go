package notes

import (
	"fmt"
	"testing"
)

// When defining a function instead of a method, be sure that it does need to be a function since a function nested inside another method or function cannot be stubbed
func Log(message string) {
	fmt.Println("Logging into file: " + message)
}

// Log function cannot be stubbed, hence, testing calculation code would need to run the actual Log function
func Calculate() {
	Log("Calculating")
	// Do calculation code here
}

// There is no way to stub Log as it is a function
func TestCalculate(t *testing.T) {
	Calculate()
}
