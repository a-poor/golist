package main

import (
	"errors"
	"time"

	"github.com/a-poor/golist"
)

func main() {
	// Create a new list
	l := golist.NewList()
	l.FailOnError = true

	// This task runs as planned
	l.AddTask(&golist.Task{
		Message: "Get a pen",
		Action: func(c golist.TaskContext) error {
			time.Sleep(time.Second)
			return nil
		},
	})

	// This task skips itself
	l.AddTask(&golist.Task{
		Message: "Get some paper",
		Skip: func(c golist.TaskContext) bool {
			return true
		},
		Action: func(c golist.TaskContext) error {
			time.Sleep(time.Second)
			return nil
		},
	})

	// This task returns an error
	l.AddTask(&golist.Task{
		Message: "Build a desk",
		Action: func(c golist.TaskContext) error {
			time.Sleep(time.Second)
			return errors.New("not good at carpentry")
		},
	})

	// These tasks are skipped because `FailOnError` is set
	// to `true` and a previous task failed
	l.AddTask(&golist.Task{
		Message: "Write a novel",
		Action: func(c golist.TaskContext) error {
			time.Sleep(time.Second)
			return nil
		},
	})
	l.AddTask(&golist.Task{
		Message: "Become famous",
		Action: func(c golist.TaskContext) error {
			time.Sleep(time.Second)
			return nil
		},
	})

	// Run the tasks
	l.RunAndWait()
}
