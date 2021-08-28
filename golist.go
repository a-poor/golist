package main

import (
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
	list := List{}

	list.AddTask(&Task{
		Message: "Task 1",
		Action: func() error {
			time.Sleep(time.Second / 2)
			return nil
		},
	})
	list.AddTask(&Task{
		Message: "Task 2",
		Action: func() error {
			time.Sleep(time.Second / 4)
			return nil
		},
	})
	tg := TaskGroup{Message: "Task Group 3"}
	tg.AddTask(&Task{
		Message: "Task 3a",
		Action: func() error {
			time.Sleep(time.Second / 4)
			return nil
		},
	})
	tg.AddTask(&Task{
		Message: "Task 3b",
		Action: func() error {
			time.Sleep(time.Second / 2)
			return nil
		},
	})
	tg.AddTask(&Task{
		Message: "Task 3c",
		Action: func() error {
			time.Sleep(time.Second / 3)
			return nil
		},
	})
	list.AddTask(&tg)
	list.AddTask(&Task{
		Message: "Task 4",
		Action: func() error {
			time.Sleep(time.Second / 3)
			return nil
		},
	})

	fmt.Println("Starting...")
	list.Start()
	list.Run()
	list.Stop()
	fmt.Println("Done.")

}
