package main

import (
	"errors"
	"fmt"
	"time"
)

var StatusMsg = map[TaskStatus]string{
	TaskNotStarted: "Not Started",
	TaskInProgress: "In Progress",
	TaskCompleted:  "Completed",
	TaskFailed:     "Failed",
	TaskSkipped:    "Skipped",
}

func main() {
	list := NewList()
	list.AddNewTask("Do this first", func() error {
		time.Sleep(time.Millisecond * 1000)
		return nil
	})
	t2 := Task{
		Title: "Skip this",
		Skip: func() bool {
			return true
		},
		Action: func() error {
			time.Sleep(time.Millisecond * 500)
			return nil
		},
	}
	list.AddTask(&t2)
	list.AddNewTask("Then try this", func() error {
		time.Sleep(time.Millisecond * 1000)
		return errors.New("oh no")
	})
	list.AddNewTask("And finally do", func() error {
		time.Sleep(time.Millisecond * 1000)
		return nil
	})
	fmt.Println("Starting...")
	list.Start()
	time.Sleep(time.Millisecond * 1000)
	list.Run()
	fmt.Println("Done.")

}
