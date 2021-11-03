package golist

import (
	"sync"

	"github.com/hashicorp/go-multierror"
)

// TaskGroup represents a group of TaskRunners
// for running nested tasks within a TaskList
type TaskGroup struct {
	Message                 string                 // The message to be displayed
	Tasks                   []TaskRunner           // A list of tasks to run
	Skip                    func(TaskContext) bool // Is run before the task starts. If returns true, the task isn't run
	FailOnError             bool                   // If true, the task group stops on the first error
	HideTasksWhenNotRunning bool                   // If true, the task group only show its sub-task-runners when actively running
	Concurrent              bool                   // Should the tasks be run concurrently?

	status TaskStatus // The status of the task

	lock sync.RWMutex
}

// NewTaskGroup creates a new TaskGroup
// with the message `m` and the tasks `ts`.
func NewTaskGroup(m string, ts []TaskRunner) *TaskGroup {
	return &TaskGroup{
		Message: m,
		Tasks:   ts,
	}
}

// AddTask adds a TaskRunner to this TaskGroup's tasks
// and returns a pointer to itself.
func (tg *TaskGroup) AddTask(t TaskRunner) *TaskGroup {
	tg.lock.Lock()
	defer tg.lock.Unlock()
	tg.Tasks = append(tg.Tasks, t)
	return tg
}

// runSync runs the TaskRunners in this TaskGroup synchronously
func (tg *TaskGroup) runSync(c TaskContext) error {
	var skipRemaining bool
	for _, t := range tg.Tasks {
		if skipRemaining {
			t.SetStatus(TaskSkipped)
			continue
		}
		err := t.Run(c)
		if err != nil && tg.FailOnError {
			skipRemaining = true
		}
	}
	return tg.GetError()
}

// runAsync runs the TaskRunners in this TaskGroup concurrently
// and blocks until they are all done.
func (tg *TaskGroup) runAsync(c TaskContext) error {
	var wg sync.WaitGroup
	for _, t := range tg.Tasks {
		wg.Add(1)
		go func(t TaskRunner) {
			defer wg.Done()
			t.Run(c)
		}(t)
	}
	wg.Wait()
	return tg.GetError()
}

// Run runs the TaskRunners in this TaskGroup
func (tg *TaskGroup) Run(parentContext TaskContext) error {
	// Create a context
	c := tg.createContext(parentContext)

	// Check if the task should be skipped
	if tg.Skip != nil && tg.Skip(c) {
		tg.SetStatus(TaskSkipped)
		return nil
	}

	// Prepare to run
	tg.SetStatus(TaskInProgress)

	var err error
	if tg.Concurrent {
		// Run concurrently...
		err = tg.runAsync(c)
	} else {
		// Otherwise, run synchronously
		err = tg.runSync(c)
	}

	// Update the TaskGroup's status
	if err != nil {
		tg.SetStatus(TaskFailed)
	} else {
		tg.SetStatus(TaskCompleted)
	}

	// Return the error
	return err
}

// createContext creates a TaskContext for the task
func (t *TaskGroup) createContext(parentContext TaskContext) TaskContext {
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

// SetMessage sets the display message for this TaskGroup
func (tg *TaskGroup) SetMessage(m string) {
	tg.lock.Lock()
	defer tg.lock.Unlock()
	tg.Message = m
}

// SetMessage sets the display message for this TaskGroup
func (tg *TaskGroup) GetMessage() string {
	tg.lock.RLock()
	defer tg.lock.RUnlock()
	return tg.Message
}

// GetError returns this TaskGroup's errors, if any
func (tg *TaskGroup) GetError() error {
	var err *multierror.Error
	for _, t := range tg.Tasks {
		err = multierror.Append(err, t.GetError())
	}
	return err.ErrorOrNil()
}

// GetStatus returns this TaskGroup's TaskStatus
func (tg *TaskGroup) GetStatus() TaskStatus {
	tg.lock.RLock()
	defer tg.lock.RUnlock()
	return tg.status
}

// GetStatus sets this TaskGroup's TaskStatus
func (tg *TaskGroup) SetStatus(s TaskStatus) {
	tg.lock.Lock()
	defer tg.lock.Unlock()
	tg.status = s
}

// GetTaskStates returns a slice of TaskStates for this TaskGroup
// representing it's current state as well as the state of its sub-tasks
// (by calling GetTaskStates on each of its sub-tasks). TaskStates store
// a TaskRunners message, status, and tree-depth, and are passed up to
// the parent List for printing.
func (tg *TaskGroup) GetTaskStates() []*TaskState {
	messages := []*TaskState{{
		Status:  tg.GetStatus(),
		Message: tg.GetMessage(),
	}}
	if !tg.HideTasksWhenNotRunning || tg.GetStatus() == TaskInProgress {
		for _, t := range tg.Tasks {
			msgs := t.GetTaskStates()
			for _, m := range msgs {
				m.Depth++
			}
			messages = append(messages, msgs...)
		}
	}
	return messages
}
