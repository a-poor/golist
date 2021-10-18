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
