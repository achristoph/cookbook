// 5. Testing conditional logic where one case does not have a value to test for, e.g: logger.Log("this case is skipped"). In go, reflection cannot help the issue either, so we can add a field that keep track of the message being logged
package notes

import (
	"fmt"
	"testing"
)

type Logger interface {
	Log(message string)
	AddMessage(message string)
	HasMessage(message string) bool
}

type ConsoleLogger struct {
	Logger
	Messages map[string]bool
}

func (c ConsoleLogger) Log(message string) {
	c.AddMessage(message)
	fmt.Println(message)
}

func (c ConsoleLogger) AddMessage(message string) {
	c.Messages[message] = true
}
func (c ConsoleLogger) HasMessage(message string) bool {
	return c.Messages[message]
}

type Calculator struct {
	SkipExtra bool
	Logger
}

func (c Calculator) ExtraRun() {
	c.Logger.Log("Running extra")
}

func (c Calculator) Run() {
	if c.SkipExtra {
		// Log does not return a value, so we need to use extra field to help with the test
		c.Logger.Log("Extra skipped")
	} else {
		c.ExtraRun()
	}
}

func TestCalculator_NormalRun(t *testing.T) {
	consoleLogger := ConsoleLogger{Messages: make(map[string]bool)}
	calculator := Calculator{SkipExtra: true, Logger: consoleLogger}
	calculator.Run()
	if !calculator.Logger.HasMessage(("Extra skipped")) {
		t.Fatal("Not found")
	}
}

func TestCalculator_SkipExtraRun(t *testing.T) {
	calculator := Calculator{SkipExtra: false, Logger: ConsoleLogger{Messages: make(map[string]bool)}}
	calculator.Run()
	if !calculator.Logger.HasMessage(("Running extra")) {
		t.Fatal("Not found")
	}
}
