package main

import (
	"time"

	"github.com/a-poor/golist"
)

func main() {
	// Create a new list
	l := golist.NewList()
	// l.StatusIndicator = nil

	// Add some tasks
	l.AddTask(golist.NewTask("Get a pen", func(c golist.TaskContext) error {
		time.Sleep(time.Second)
		// time.Sleep(time.Millisecond * 100)
		return nil
	}))
	l.AddTask(golist.NewTask("Get some paper", func(c golist.TaskContext) error {
		time.Sleep(time.Second)
		// time.Sleep(time.Millisecond * 100)
		return nil
	}))
	l.AddTask(golist.NewTask("Write a novel", func(c golist.TaskContext) error {
		time.Sleep(time.Second)
		// time.Sleep(time.Millisecond * 100)
		return nil
	}))

	l.Start()
	time.Sleep(time.Millisecond * 100)

	// Run the tasks
	l.Run()
	time.Sleep(time.Millisecond * 100)

	// Stop dispalying the list status
	l.Stop()
	// fmt.Println("Almost Done.")
	time.Sleep(time.Millisecond * 100)
	// fmt.Println("Done.")
	// fmt.Println("Done.")
	// fmt.Println("Done.")
}
