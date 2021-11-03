package golist

import "sync"

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
	Skip    func(TaskContext) bool  // Is run before the task starts. If returns true, the task isn't run

	status TaskStatus // The status of the task
	err    error      // The error returned by the task function

	lock sync.RWMutex
}

// NewTask creates a new Task with the message `m`
// and the action function `a`.
func NewTask(m string, a func(TaskContext) error) *Task {
	return &Task{
		Message: m,
		Action:  a,
	}
}

// Run runs the task's action function
func (t *Task) Run(parentContext TaskContext) error {
	// Create a TaskContext to pass to `Skip` and `Action`
	c := t.createContext(parentContext)

	// Check if the task should be skipped
	if t.Skip != nil && t.Skip(c) {
		t.SetStatus(TaskSkipped)
		return nil
	}

	// Check that an action function exists
	if t.Action == nil {
		t.SetStatus(TaskFailed)
		t.status = TaskFailed
		t.SetError(ErrNilAction)
		return ErrNilAction
	}

	// Set the status to in-progress and run
	t.SetStatus(TaskInProgress)
	err := t.Action(c)

	// Evaluate the error and update the task status
	if err != nil {
		t.SetStatus(TaskFailed)
	} else {
		t.SetStatus(TaskCompleted)
	}

	// Store the error and return it
	t.SetError(err)
	return err
}

// createContext creates a TaskContext for the task
func (t *Task) createContext(parentContext TaskContext) TaskContext {
	return &taskContext{
		setMessage: func(msg string) {
			t.SetMessage(msg)
		},
		println: func(a ...interface{}) error {
			return parentContext.Println(a...)
		},
		printfln: func(f string, a ...interface{}) error {
			return parentContext.Printfln(f, a...)
		},
	}
}

// SetMessage sets the Task's message text
func (t *Task) SetMessage(m string) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.Message = m
}

// GetMessage returns Task's message text
func (t *Task) GetMessage() string {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.Message
}

// SetError sets the Task's error value
func (t *Task) SetError(err error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.err = err
}

// GetError returns Task's error value, if there is one
func (t *Task) GetError() error {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.err
}

// GetStatus returns the Task's status
func (t *Task) GetStatus() TaskStatus {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.status
}

// SetStatus sets the Task's status
func (t *Task) SetStatus(s TaskStatus) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.status = s
}

// GetTaskTates returns the TaskState description
// of the current task
func (t *Task) GetTaskStates() []*TaskState {
	return []*TaskState{{
		Message: t.GetMessage(),
		Status:  t.GetStatus(),
	}}
}
