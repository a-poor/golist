package golist

import (
	"errors"
	"testing"
)

func TestNewTaskGroup(t *testing.T) {
	group := NewTaskGroup("test", []TaskRunner{
		NewTask("t0", func(c TaskContext) error {
			return nil
		}),
		NewTask("t1", func(c TaskContext) error {
			return nil
		}),
	})

	err := group.Run(&taskContext{})

	if err != nil {
		t.Fatal(err)
	}
}

func TestNewTaskGroupWithErr(t *testing.T) {
	group := NewTaskGroup("test", []TaskRunner{
		NewTask("t0", func(c TaskContext) error {
			return nil
		}),
		NewTask("t1", func(c TaskContext) error {
			return errors.New("task error")
		}),
	})

	err := group.Run(&taskContext{})

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestTaskGroup_GetTaskStates(t *testing.T) {
	g := NewTaskGroup("test", []TaskRunner{
		NewTask("t0", func(c TaskContext) error {
			return nil
		}),
	})
	g.AddTask(NewTask("t1", func(c TaskContext) error {
		return nil
	}))

	s := g.GetTaskStates()
	for _, k := range s {
		if k.Status != TaskNotStarted {
			t.Errorf("expected task state to be %q, got %q", TaskNotStarted, k.Status)
		}
	}
}

func TestTaskGroup_GetStatus(t *testing.T) {
	g := NewTaskGroup("test", []TaskRunner{})
	if s := g.GetStatus(); s != TaskNotStarted {
		t.Errorf("expected task status to be %q, got %q", TaskNotStarted, s)
	}
}

func TestTaskGroup_SetMessage(t *testing.T) {
	start := "start"
	end := "end"

	g := NewTaskGroup(start, []TaskRunner{})
	g.SetMessage(end)
	m := g.Message

	if m != end {
		t.Errorf("expected message to be %q, got %q", end, m)
	}
}

func TestTaskGroup_createContext(t *testing.T) {
	var setMsgCalled bool
	var printlnCalled bool
	var printfCalled bool

	parent := &taskContext{
		setMessage: func(msg string) {
			setMsgCalled = true
		},
		println: func(a ...interface{}) error {
			printlnCalled = true
			return nil
		},
		printfln: func(msg string, a ...interface{}) error {
			printfCalled = true
			return nil
		},
	}

	old := "old"
	new := "new"
	g := NewTaskGroup(old, []TaskRunner{})
	c := g.createContext(parent)

	c.SetMessage(new)
	c.Println()
	c.Printfln("")

	if setMsgCalled {
		t.Error("expected SetMessage to replaced")
	}
	if m := g.Message; m != new {
		t.Errorf("expected message to be %q, got %q", new, m)
	}

	if !printlnCalled {
		t.Error("expected Println to be called")
	}
	if !printfCalled {
		t.Error("expected Printfln to be called")
	}
}
