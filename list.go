package golist

import (
	"context"
	"errors"
	"io"
	"os"
	"time"
)

// Default List print delay
var DefaultListDelay = time.Millisecond * 100

// List is the top-level task list.
type List struct {
	Tasks           []TaskRunner     // List of tasks to run
	Delay           time.Duration    // Delay between prints
	FailOnError     bool             // If true, the task execution stops on the first error
	Concurrent      bool             // Should the tasks be run concurrently? NOTE: Not supported yet
	Writer          io.Writer        // Writer to use for printing output
	MaxLineLength   int              // Maximum line length for printing (0 = no limit)
	StatusIniicator StatusIndicators // Map of statuses to status indicators

	running bool               // Is the list running?
	cancel  context.CancelFunc // A context cancel function for stopping the list run
}

// NewList creates a new task list with some sensible defaults.
// It writes to stdout and and has a delay of 100ms between prints.
func NewDefaultList() *List {
	return &List{
		Writer:          os.Stdout,
		Delay:           DefaultListDelay,
		StatusIniicator: CreateDefaultStatusIndicator(),
	}
}

// NewList creates a new task list that writes to the
// provided io.Writer. Mostly used for testing.
func NewListWithWriter(w io.Writer) *List {
	return &List{
		Writer:          w,
		Delay:           DefaultListDelay,
		StatusIniicator: CreateDefaultStatusIndicator(),
	}
}

// AddTask adds a TaskRunner to the top-level List
func (l *List) AddTask(t TaskRunner) {
	if l.Tasks == nil {
		l.Tasks = make([]TaskRunner, 0)
	}
	l.Tasks = append(l.Tasks, t)
}

// Start begins displaying the list statuses
// from a background goroutine.
//
// Note: If the list is created without a writer,
// this function will return an error.
func (l *List) Start() error {
	if l.Writer == nil {
		return errors.New("no writer specified")
	}

	// Check if it's already displaying
	if l.running {
		return nil
	}

	// Create a cancelable context
	ctx, cancel := context.WithCancel(context.Background())
	l.cancel = cancel

	// Start the display loop
	go func() {
		l.print()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				l.clear()
				l.print()
				time.Sleep(time.Millisecond * 100)
			}
		}
	}()

	// Set the running flag
	l.running = true

	return nil
}

// Run starts calling the Action functions for each task
func (l *List) Run() error {
	l.Start()
	defer l.Stop()

	for _, t := range l.Tasks {
		if err := t.Run(); err != nil && l.FailOnError {
			return err
		}
	}

	return nil
}

// Stop stops displaying the task list statuses and
// cancels the background goroutine.
//
// Note: Stop also calls the `clear` and `print` functions
// one final time each before finishing.
func (l *List) Stop() {
	// Check if it's already NOT displaying
	if !l.running {
		return
	}

	// Send the cancel signal
	l.cancel()

	// Clear and print one final time (NOTE: should this be an option?)
	l.clear()
	l.print()

	l.running = false
	l.cancel = nil
}

// print calls the `print` function for each of the TaskRunners
func (l *List) print() {
	for _, t := range l.Tasks {
		t.Print(0)
	}
}

// clear calls the `clear` function for each of the TaskRunners
func (l *List) clear() {
	for _, t := range l.Tasks {
		t.Clear()
	}
}

func (l *List) getTaskStates() []TaskState {
	var messages []TaskState
	for _, t := range l.Tasks {
		msgs := t.GetTaskStates()
		for _, m := range msgs {
			m.Depth++
		}
		messages = append(messages, msgs...)
	}
	return messages
}

// type TaskState struct {
// 	Message string
// 	Status  TaskStatus
// 	Depth   int
// }

func (l *List) formatMessage(m TaskState) string {
	return ""
}

func (l *List) print2(states []TaskState) {

}
