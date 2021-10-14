package golist

import (
	"context"
	"time"
)

type List struct {
	Tasks       []TaskRunner // List of tasks to run
	FailOnError bool         // If true, the task execution stops on the first error

	running bool               // Is the list running?
	cancel  context.CancelFunc // A context cancel function for stopping the list run
}

// NewList creates a new task list
func NewList() *List {
	return &List{}
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
func (l *List) Start() {
	// Check if it's already displaying
	if l.running {
		return
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
