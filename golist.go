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
	list.AddNewTask("Then do this", func() error {
		time.Sleep(time.Millisecond * 500)
		return nil
	})
	list.AddNewTask("And finally this", func() error {
		time.Sleep(time.Millisecond * 1000)
		return errors.New("oh no")
	})

	fmt.Println("Starting...")
	list.Run()
	fmt.Println("Done.")

}
