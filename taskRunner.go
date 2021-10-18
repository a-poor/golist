package golist

type TaskRunner interface {
	Run(TaskContext) error       // Run the task and return any error
	SetMessage(string)           // Set the task message
	GetStatus() TaskStatus       // Get the task status
	SetStatus(TaskStatus)        // Set the task status
	GetError() error             // Set the task error
	GetTaskStates() []*TaskState // Get the task states
}
