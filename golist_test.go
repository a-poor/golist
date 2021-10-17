package golist_test

import (
	"errors"
	"github.com/a-poor/golist"
	"testing"
)

func TestGolist(t *testing.T) {
	golist.NewDefaultList()
}

func TestGolist_Run(t *testing.T) {
	list := golist.NewDefaultList()

	list.AddTask(&golist.Task{
		Message: "t0",
		Action: func(c golist.TaskContext) error {
			return nil
		},
	})

	list.AddTask(&golist.Task{
		Message: "t1",
		Action: func(c golist.TaskContext) error {
			return nil
		},
	})

	// Start displaying the task status
	if err := list.Start(); err != nil {
		t.Fatal(err)
	}

	// Run the tasks
	if err := list.Run(); err != nil {
		t.Fatal(err)
	}

	// Stop displaying the task status
	list.Stop()
}

func TestGolist_RunWithErrors(t *testing.T) {
	list := golist.NewDefaultList()

	list.AddTask(&golist.Task{
		Message: "t0",
		Action: func(c golist.TaskContext) error {
			return nil
		},
	})

	list.AddTask(&golist.Task{
		Message: "t1",
		Action: func(c golist.TaskContext) error {
			return errors.New("oops")
		},
	})

	// Start displaying the task status
	if err := list.Start(); err != nil {
		t.Fatal(err)
	}

	// Run the tasks
	if err := list.Run(); err != nil {
		t.Fatal(err)
	}

	// Stop displaying the task status
	list.Stop()
}
