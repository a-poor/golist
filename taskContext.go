package golist

// TaskContext is the context passed to the Tasks'
// Action and Skip functions.
type TaskContext interface {
	SetMessage(string)                     // Set the task's message
	Println(...interface{}) error          // Safely print between list updates like `fmt.Println`
	Printfln(string, ...interface{}) error // Safely print formatted text between list updates like `fmt.Printf` but with a newline character at the end
}

// taskContext implements the TaskContext interface for
// being passed to a Task's Action and Skip functions.
type taskContext struct {
	setMessage func(string)
	println    func(...interface{}) error
	printfln   func(string, ...interface{}) error
}

// SetMessage updates the task's status message
// while running
func (tc *taskContext) SetMessage(msg string) {
	tc.setMessage(msg)
}

// Println prints text safely between list updates
func (tc *taskContext) Println(a ...interface{}) error {
	return tc.println(a...)
}

// Printfln prints formatted text safely between list updates
func (tc *taskContext) Printfln(f string, a ...interface{}) error {
	return tc.printfln(f, a...)
}
