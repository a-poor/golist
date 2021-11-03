package golist

// TaskRunner represents a task or group
// of tasks to be run.
type TaskRunner interface {
	Run(TaskContext) error       // Run the task and return any error
	GetMessage() string          // Get the task's message
	SetMessage(string)           // Set the task's message
	GetStatus() TaskStatus       // Get the task's status
	SetStatus(TaskStatus)        // Set the task's status
	GetError() error             // Set the task's error value
	GetTaskStates() []*TaskState // Return the task or subtask states and display information
}
