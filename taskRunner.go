package main

const IndentSize = 2

type TaskRunner interface {
	Run() error // Run the task and return any error
	// String(int)           // Print the task message and status indicator
	SetMessage(string)    // Set the task message
	SetStatus(TaskStatus) // Set the task status
	SetError(error)       // Set the task error
	GetSize() int         // Return the task size (number of tasks and subtasks for clearing)
	Print(int)            // Print the task message and status indicator (with an indent according to the specified depth)
	Clear()               // Clear the output of the task
}
