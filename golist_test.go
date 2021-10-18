package golist

import (
	"bytes"
	"errors"
	"os"
	"testing"
)

func TestGolist(t *testing.T) {
	NewList()
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

func TestStatus(t *testing.T) {
	if s := TaskNotStarted.String(); s != "Not Started" {
		t.Errorf("TaskNotStarted.String = %q", s)
	}
	if s := TaskInProgress.String(); s != "In Progress" {
		t.Errorf("TaskInProgress.String = %q", s)
	}
	if s := TaskCompleted.String(); s != "Completed" {
		t.Errorf("TaskCompleted.String = %q", s)
	}
	if s := TaskFailed.String(); s != "Failed" {
		t.Errorf("TaskFailed.String = %q", s)
	}
	if s := TaskSkipped.String(); s != "Skipped" {
		t.Errorf("TaskSkipped.String = %q", s)
	}

	other := TaskStatus(999)
	e := "Unknown"
	if s := other.String(); s != e {
		t.Errorf("unexpected task status expected %q, got %q", e, s)
	}
}

func TestList_NoWriter(t *testing.T) {
	l := List{}
	l.Start()
	if l.Writer == nil {
		t.Error("list's writer should have auto-set")
	} else if l.Writer != os.Stdout {
		t.Errorf("list's writer auto-set to something other than stdout, %q", l.Writer)
	}
}
