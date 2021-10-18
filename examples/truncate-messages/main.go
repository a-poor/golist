package main

import (
	"time"

	"github.com/a-poor/golist"
)

func main() {
	// Create a new list
	l := golist.NewList()
	l.MaxLineLength = 30

	// Add some tasks
	l.AddTask(golist.NewTask("Here's a really long message that should be truncated", func(c golist.TaskContext) error {
		time.Sleep(time.Second)
		return nil
	}))
	l.AddTask(golist.NewTask("This line isn't too long", func(c golist.TaskContext) error {
		time.Sleep(time.Second)
		return nil
	}))
	l.AddTask(golist.NewTask("And this line also has too much text so it should be trimmed", func(c golist.TaskContext) error {
		time.Sleep(time.Second)
		return nil
	}))

	// Run the tasks
	l.Run()

	// Stop dispalying the list status
	l.Stop()
}
