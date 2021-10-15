package golist

// TaskContext
type TaskContext interface {
	SetMessage(string) // Set the task's message

	// Copy() TaskContext // Copy the task context
	// SetKey(string, interface{}) // Store a key/value pair that's accessable by other tasks
	// GetKey(string) interface{}  // Get a key/value pair
	// DelKey(string)              // Delete a key/value pair
}

type taskContext struct {
	setMessage func(string)
}

func (tc *taskContext) SetMessage(msg string) {
	tc.setMessage(msg)
}
