package notes

import (
	"fmt"
	"testing"
)

// - When defining a function / method, there are three types of parameters:
// - literal values, e.g: title
// - pure struct: a pure struct has not method, hence, only its fields need to be changed during testing, e.g: Person
// - struct with methods: this struct has methods, hence, always define an interface and use this interface instead so that the struct methods can be stubbed during testing
type Person struct {
	Name string
	Age  int
}
type Ball struct {
	Diameter int
}

func (b Ball) Move() {
	fmt.Println("Ball is moving")
}

func (b Ball) Bounce() {
	fmt.Println("Ball is bouncing")
}

type Thing interface {
	Move()
	Bounce()
}

func FnWithActualStruct(title string, person Person, ball Ball) (output string) {
	ball.Move()
	output = fmt.Sprintf("Title: %s, Name: %s, Age: %d", title, person.Name, person.Age)
	return output
}

func FnWithInterface(title string, person Person, thing Thing) (output string) {
	thing.Move()
	output = fmt.Sprintf("Title: %s, Name: %s, Age: %d", title, person.Name, person.Age)
	return output
}

func TestFnNoStub(t *testing.T) {
	output := FnWithActualStruct("subject", Person{Name: "john", Age: 1}, Ball{Diameter: 1})
	t.Log(output)
}

type BallStub struct {
	Ball
}

// The Move is overridden
func (b BallStub) Move() {
	fmt.Println("Stubbed Move")
}

func TestFnWithStub(t *testing.T) {
	// Ball cannot be stubbed
	output := FnWithActualStruct("subject", Person{Name: "john", Age: 1}, Ball{Diameter: 1})
	t.Log(output)
	// Using BallStub that fulfills Thing interface
	output = FnWithInterface("subject", Person{Name: "john", Age: 1}, BallStub{})
	t.Log(output)
}
