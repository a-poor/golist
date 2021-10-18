package main

import (
	"time"

	"github.com/a-poor/golist"
)

func doSomething(c golist.TaskContext) error {
	time.Sleep(time.Second)
	return nil
}

func main() {
	// Create a new list
	l := golist.NewList()

	// Add some tasks
	l.AddTask(golist.NewTask("Get a pen", doSomething))

	// Add nested tasks
	l.AddTask(golist.NewTaskGroup("Get some paper", []golist.TaskRunner{

		// This TaskGroup will hide it's subtasks when they aren't running
		&golist.TaskGroup{
			Message:                 "Drive to the store",
			HideTasksWhenNotRunning: true,
			Tasks: []golist.TaskRunner{
				golist.NewTask("Leave the driveway", doSomething),
				golist.NewTask("Make some turns", doSomething),
				golist.NewTask("Enter the store's parking lot", doSomething),
			},
		},
		golist.NewTask("Get a box of paper", doSomething),
	}))

	// This TaskGroup will run it's subtasks concurrently
	l.AddTask(&golist.TaskGroup{
		Message:    "Write a novel",
		Concurrent: true,
		Tasks: []golist.TaskRunner{
			golist.NewTask("Write the first chapter", doSomething),
			golist.NewTask("Write the second chapter", doSomething),
			golist.NewTask("Write the third chapter", doSomething),
		},
	})
	l.AddTask(golist.NewTask("Publish the novel", doSomething))

	// Run the tasks
	l.RunAndWait()
}
