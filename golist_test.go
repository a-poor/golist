package golist

import (
	"bytes"
	"errors"
	"testing"
)

func TestGolist(t *testing.T) {
	_ = NewList()
}

func TestGolist_Run(t *testing.T) {
	list := NewList()
	list.Writer = &bytes.Buffer{}

	list.AddTask(&Task{
		Message: "t0",
		Action: func(c TaskContext) error {
			return nil
		},
	})

	list.AddTask(&Task{
		Message: "t1",
		Action: func(c TaskContext) error {
			return nil
		},
	})

	// Start displaying the task status
	list.Start()

	// Run the tasks
	if err := list.Run(); err != nil {
		t.Fatal(err)
	}

	// Stop displaying the task status
	list.Stop()
}

func TestGolist_RunWithErrors(t *testing.T) {
	list := NewList()
	list.Writer = &bytes.Buffer{}

	list.AddTask(&Task{
		Message: "t0",
		Action: func(c TaskContext) error {
			return nil
		},
	})

	list.AddTask(&Task{
		Message: "t1",
		Action: func(c TaskContext) error {
			return errors.New("oops")
		},
	})

	// Start displaying the task status
	list.Start()

	// Run the tasks
	if err := list.Run(); err == nil {
		t.Fatalf("Expected 1 error and got %q", err)
	}

	// Stop displaying the task status
	list.Stop()
}
