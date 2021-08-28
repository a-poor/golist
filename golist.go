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
	list.AddTask(&Task{
		Message: "Task 3",
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

	// si := createStatusIndicator(statusIndicatorConfig{})
	// fmt.Println(si(TaskNotStarted))
	// time.Sleep(time.Millisecond * 100)

	// fmt.Println(si(TaskInProgress))
	// time.Sleep(time.Millisecond * 100)
	// fmt.Println(si(TaskInProgress))
	// time.Sleep(time.Millisecond * 100)
	// fmt.Println(si(TaskInProgress))
	// time.Sleep(time.Millisecond * 100)
	// fmt.Println(si(TaskInProgress))
	// time.Sleep(time.Millisecond * 100)

	// fmt.Println(si(TaskCompleted))
	// time.Sleep(time.Millisecond * 100)
	// fmt.Println(si(TaskFailed))
	// time.Sleep(time.Millisecond * 100)
	// fmt.Println(si(TaskSkipped))
	// time.Sleep(time.Millisecond * 100)

}
