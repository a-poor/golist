package golist

import (
	"fmt"
	"strings"
)

type TaskState struct {
	Message string
	Status  TaskStatus
	Depth   int
}

// Task represents a task to be run as part
// of a List or TaskGroup
type Task struct {
	Message string                  // Message to display to user
	Action  func(TaskContext) error // The task function to be run
	Skip    func() bool             // Is run before the task starts. If returns true, the task isn't run

	statusIndicator func(TaskStatus) string
	status          TaskStatus // The status of the task
	err             error      // The error returned by the task function
}

// Run runs the task's action function
func (t *Task) Run() error {
	if t.Skip != nil && t.Skip() {
		t.status = TaskSkipped
		return nil
	}

	if t.Action == nil {
		t.status = TaskFailed
		t.err = ErrNilAction
		return t.err
	}

	t.status = TaskInProgress
	ctx := t.createContext()
	err := t.Action(ctx)

	if err != nil {
		t.status = TaskFailed
	} else {
		t.status = TaskCompleted
	}
	t.err = err
	return err
}

func (t *Task) init() {
	t.initStatusIndicator()
}

// createContext creates a TaskContext for the task
func (t *Task) createContext() TaskContext {
	return &taskContext{func(msg string) {
		t.SetMessage(msg)
	}}
}

func (t *Task) Print(indent int) {
	if t.statusIndicator == nil {
		t.init()
	}
	pad := strings.Repeat(" ", indent)
	stat := t.statusIndicator(t.status)
	fmt.Printf("%s%s %s\n", pad, stat, t.Message)
}

func (t *Task) Clear() {
	fmt.Print("\033[1A") // Move up a line
	fmt.Print("\033[K")  // Clear the line
	fmt.Print("\r")      // Move back to the beginning of the line
}

// SetMessage sets the Task's message text
func (t *Task) SetMessage(m string) {
	t.Message = m
}

// SetError sets the Task's error value
func (t *Task) SetError(err error) {
	t.err = err
}

// GetStatus returns the Task's status
func (t *Task) GetStatus() TaskStatus {
	return t.status
}

// SetStatus sets the Task's status
func (t *Task) SetStatus(s TaskStatus) {
	t.status = s
}

func (t *Task) GetDepth() int {
	return 1
}

func (t Task) GetSize() int {
	return 1
}

func (t *Task) initStatusIndicator() {
	t.statusIndicator = createStatusIndicator(statusIndicatorConfig{})
}

func (t *Task) GetTaskStates() []TaskState {
	return []TaskState{{
		Message: t.Message,
		Status:  t.status,
	}}
}
