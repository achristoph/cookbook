// 5. Testing conditional logic where one case does not have a value to test for (method without return value), e.g: logger.Log("this case is skipped").
// In go, reflection cannot help the issue either, so the green path can be let go without being tested - we don't know if a success path actually called the Log method

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
}

func (c ConsoleLogger) Log(message string) {
	fmt.Println(message)
}

func (c ConsoleLogger) HasMessage(message string) {
	// console logger can have a no-op
}

type Calculator struct {
	SkipExtra bool
	Logger
	Messages map[string]bool
}

func (c Calculator) ExtraRun() {
}

func (c Calculator) Run() {
	if c.SkipExtra {
		// Log does not return a value, so we need to use extra field to help with the test
		c.Logger.Log("Extra skipped")
	} else {
		c.ExtraRun()
		c.Logger.Log("Running extra")
	}
}

func TestCalculator_NormalRun(t *testing.T) {
	calculator := Calculator{SkipExtra: false, Logger: ConsoleLogger{}}
	calculator.Run()
}

func TestCalculator_SkipExtraRun(t *testing.T) {
	calculator := Calculator{SkipExtra: true, Logger: ConsoleLogger{}}
	calculator.Run()
}
