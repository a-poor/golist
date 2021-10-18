package golist

import (
	"bytes"
	"errors"
	"os"
	"testing"
	"time"
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
	list.Start()
	list.Start()

	// Run the tasks
	if err := list.Run(); err != nil {
		t.Fatal(err)
	}

	// Stop displaying the task status
	list.Stop()
	list.Stop()
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

func TestNewListWithWriter(t *testing.T) {
	w := &bytes.Buffer{}
	l := NewListWithWriter(w)
	if l.Writer != w {
		t.Errorf("list's writer should have been set to %q, got %q", w, l.Writer)
	}
}

func TestList_runAsync(t *testing.T) {
	l := NewList()
	l.Concurrent = true
	l.Writer = &bytes.Buffer{}

	var t0Start bool
	var t0Stop bool

	l.AddTask(&Task{
		Message: "t0",
		Action: func(c TaskContext) error {
			t0Start = true
			time.Sleep(time.Millisecond * 100)
			t0Stop = true
			return nil
		},
	})

	l.AddTask(&Task{
		Message: "t1",
		Action: func(c TaskContext) error {
			time.Sleep(time.Millisecond * 10)
			if !t0Start {
				t.Error("t0 should have started already")
			}
			if t0Stop {
				t.Error("t0 shouldn't have finished yet")
			}
			return nil
		},
	})

	// Start displaying the task status
	l.RunAndWait()
}

func TestList_createRootContext(t *testing.T) {
	l := NewList()
	l.Writer = &bytes.Buffer{}
	c := l.createRootContext()
	c.SetMessage("")
	c.Println("")
	c.Printfln("")
}

func TestListReturnErrors(t *testing.T) {
	l := NewList()
	l.Writer = &bytes.Buffer{}

	expect := errors.New("oh no")
	l.AddTask(NewTask("t0", func(c TaskContext) error {
		return expect
	}))

	got := l.RunAndWait()
	if errors.Unwrap(got) != expect {
		t.Errorf("expected error %q, got %q", expect, got)
	}
}

func TestListPrints(t *testing.T) {
	l := NewList()
	l.Writer = &bytes.Buffer{}
	l.AddTask(NewTask("t0", func(c TaskContext) error {
		c.Println("hello")
		c.Printfln("number: %d", 123)
		return nil
	}))
	l.RunAndWait()
}

func TestListSkipRemaining(t *testing.T) {
	l := NewList()
	l.Writer = &bytes.Buffer{}
	l.FailOnError = true

	var t0Ran bool
	var t1Ran bool
	var t2Ran bool

	l.AddTask(NewTask("t0", func(c TaskContext) error {
		t0Ran = true
		return nil
	}))
	l.AddTask(NewTask("t1", func(c TaskContext) error {
		t1Ran = true
		return errors.New("oh no")
	}))
	l.AddTask(NewTask("t2", func(c TaskContext) error {
		t2Ran = true
		return nil
	}))
	l.Start()
	l.runSync(l.createRootContext())

	if !t0Ran {
		t.Error("t0 should have run")
	}
	if !t1Ran {
		t.Error("t1 should have run")
	}
	if t2Ran {
		t.Error("t2 should have been skipped")
	}
}

func TestListTruncate(t *testing.T) {
	l := NewList()
	l.Writer = &bytes.Buffer{}

	i := 25
	m := "this should be cutoff because it's too long"
	e := "this should be cutoff be…"
	if n := l.truncateMessage(m, i); n != e {
		t.Errorf("expected %q, got %q", e, n)
	}

	i = 1
	m = "this should be cutoff because it's too long"
	e = "…"
	if n := l.truncateMessage(m, i); n != e {
		t.Errorf("expected %q, got %q", e, n)
	}

	i = 0
	m = "this should be cutoff because it's too long"
	e = "…"
	if n := l.truncateMessage(m, i); n != e {
		t.Errorf("expected %q, got %q", e, n)
	}
}

func TestListTruncate2(t *testing.T) {
	l := NewList()
	l.Writer = &bytes.Buffer{}
	l.MaxLineLength = 1
	l.AddTask(NewTask("Nothing should show here", func(c TaskContext) error {
		return nil
	}))
	l.RunAndWait()
}
