package golist

import (
	"errors"
	"testing"
)

func TestNewTask_RunComplete(t *testing.T) {
	msg := "test"
	task := NewTask(msg, func(c TaskContext) error {
		return nil
	})

	err := task.Run(&taskContext{})

	if err != nil {
		t.Fatal(err)
	}

	expectedStatus := TaskCompleted
	if task.GetStatus() != expectedStatus {
		t.Fatalf("expected status %s but was %s", expectedStatus, task.GetStatus())
	}

	if task.err != nil {
		t.Fatalf("expected task err to be <nil> but was: %s", task.err)
	}

	if task.Message != msg {
		t.Fatalf("expected task message to be %s but was %s", msg, task.Message)
	}
}

func TestNewTask_RunFailed(t *testing.T) {
	taskErr := errors.New("task error")
	task := NewTask("test", func(c TaskContext) error {
		return taskErr
	})

	err := task.Run(&taskContext{})

	if err == nil {
		t.Fatal("expected error for failed task run")
	}

	expectedStatus := TaskFailed
	if task.GetStatus() != expectedStatus {
		t.Fatalf("expected status %s but was %s", expectedStatus, task.GetStatus())
	}

	if task.err == nil {
		t.Fatalf("expected task err to be present but was nil")
	}
}

func TestNewTask_RunNotStarted(t *testing.T) {
	task := NewTask("test", func(c TaskContext) error {
		return nil
	})

	expectedStatus := TaskNotStarted
	if task.GetStatus() != expectedStatus {
		t.Fatalf("expected status %s but was %s", expectedStatus, task.GetStatus())
	}

	if task.err != nil {
		t.Fatalf("expected task err to be <nil> but was: %s", task.err)
	}
}

func TestTask_SkipRun(t *testing.T) {
	k := &Task{
		Message: "test",
		Action: func(c TaskContext) error {
			t.Error("task wasn't successfully skipped")
			return nil
		},
		Skip: func(c TaskContext) bool {
			return true
		},
	}
	k.Run(&taskContext{})
}

func TestTask_NilAction(t *testing.T) {
	k := &Task{
		Message: "test",
		Action:  nil,
	}
	err := k.Run(&taskContext{})
	if err == nil {
		t.Error("expected nil action error")
	}
	if err != nil && err != ErrNilAction {
		t.Errorf("expected err = %q, got %q", ErrNilAction, err)
	}
}

func TestTask_SetMessage(t *testing.T) {
	orgMsg := "original"
	newMsg := "new"
	k := &Task{
		Message: orgMsg,
		Action:  nil,
	}
	k.SetMessage(newMsg)
	if k.Message == orgMsg {
		t.Error("task message unchanged")
		return
	}
	if k.Message != newMsg {
		t.Errorf("task message expected %q, got %q", newMsg, k.Message)
	}
}

func TestTask_createContext(t *testing.T) {

}
