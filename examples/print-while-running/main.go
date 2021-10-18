package main

import (
	"time"

	"github.com/a-poor/golist"
)

func main() {
	// Create a new list
	l := golist.NewList()

	// Add some tasks
	l.AddTask(golist.NewTask("Get a pen", func(c golist.TaskContext) error {
		time.Sleep(time.Second)
		c.Println("Printing a message from inside a task!")
		return nil
	}))
	l.AddTask(golist.NewTask("Get some paper", func(c golist.TaskContext) error {
		time.Sleep(time.Second)
		c.Printfln("Printing a %q from inside a task!", "formatted message")
		return nil
	}))
	l.AddTask(golist.NewTask("Write a novel", func(c golist.TaskContext) error {
		c.Println("Aaaand two last messages!")
		c.Println("Done.")
		time.Sleep(time.Second)
		return nil
	}))

	// Run the tasks
	l.Run()

	// Stop dispalying the list status
	l.Stop()
}
