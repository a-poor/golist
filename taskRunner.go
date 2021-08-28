package main

const IndentSize = 2

type TaskRunner interface {
	Run() error            // Run the task and return any error
	SetMessage(string)     // Set the task message
	GetStatus() TaskStatus // Get the task status
	SetStatus(TaskStatus)  // Set the task status
	SetError(error)        // Set the task error
	GetSize() int          // Return the task size (number of tasks and subtasks for clearing)
	Print(int)             // Print the task message and status indicator (based on the specified depth)
	Clear()                // Clear the output of the task
}
