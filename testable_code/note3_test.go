package notes

import (
	"fmt"
	"testing"
)

// Any variable that you want to test inside function/method will need to be returned so that the its value can be tested, however, there are variables that are 'in the middle' that is produced by some logic that should be tested - this logic needs to be extracted out as a method so that it can be tested in isolation.
// If there a few possible values to test, put the possible inputs in a collection to iterate so that it can be done efficiently in one test method

func ColorToAction(color string) (output string) {
	var action string
	switch color {
	case "green":
		action = "go"
	case "red":
		action = "stop"
	}
	output = fmt.Sprintf("%s is %s", color, action)
	return output
}

func ColorToActionAfterExtractMethod(color string) (output string) {
	// FindAction is declared as a function instead of a method, thus, cannot be stubbed, in this context is acceptable as it does not do external service call / intensive computation or other logic that are too complicated to test
	// FindAction makes action value can be tested since the logic is now encapsulated in the function
	action := FindAction(color)
	output = fmt.Sprintf("%s is %s", color, action)
	return output
}
func FindAction(color string) (action string) {
	switch color {
	case "green":
		action = "go"
	case "red":
		action = "stop"
	}
	return action
}

func TestFindAction(t *testing.T) {
	color := "green"
	actualAction := FindAction(color)
	expectedAction := "go"
	if actualAction != expectedAction {
		t.Errorf("%s not equal to %s", actualAction, expectedAction)
	}
}

// If there a few possible values to test, put the possible inputs in a collection to iterate
// so that it can be done efficiently in one test method
func TestColorToAction_CheckAllInputs(t *testing.T) {
	colors := []string{"red", "green"}
	for _, c := range colors {
		actualOutput := ColorToActionAfterExtractMethod(c)
		expectedOutput := "red is stop"
		if actualOutput != expectedOutput {
			t.Errorf("%s not equal to %s", actualOutput, expectedOutput)
		}
	}
}
