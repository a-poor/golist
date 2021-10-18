package main

import (
	"time"

	"github.com/a-poor/golist"
)

func main() {
	// Create a new list
	l := golist.NewList()

	// Add some tasks that update their messages while running
	l.AddTask(golist.NewTask("Get a pen", func(c golist.TaskContext) error {
		c.SetMessage("Getting a pen...")
		time.Sleep(time.Second)
		c.SetMessage("Got a pen!")
		return nil
	}))
	l.AddTask(golist.NewTask("Get some paper", func(c golist.TaskContext) error {
		c.SetMessage("Getting some paper...")
		time.Sleep(time.Second)
		c.SetMessage("Got some paper!")
		return nil
	}))
	l.AddTask(golist.NewTask("Write a novel", func(c golist.TaskContext) error {
		c.SetMessage("Writing a novel...")
		time.Sleep(time.Second)
		c.SetMessage("Wrote a novel!")
		return nil
	}))

	// Start the list
	l.Start()

	// Wait a second before starting the tasks to see the initial messages
	time.Sleep(time.Second)

	// Start the tasks
	l.RunAndWait()
}
