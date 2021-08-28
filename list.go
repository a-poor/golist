package main

import (
	"context"
	"time"
)

type List struct {
	Tasks       []TaskRunner  // List of tasks to run
	FailOnError bool          // If true, the task execution stops on the first error
	Delay       time.Duration // Delay between each task execution

	running bool
	cancel  context.CancelFunc
}

func NewList() *List {
	return &List{}
}

func (l *List) AddTask(t TaskRunner) {
	if l.Tasks == nil {
		l.Tasks = make([]TaskRunner, 0)
	}
	l.Tasks = append(l.Tasks, t)
}

// Start to display the list status
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

// Start to run the task list
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

// Stop displaying the task list
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

func (l *List) print() {
	for _, t := range l.Tasks {
		t.Print(0)
	}
}

func (l *List) clear() {
	for _, t := range l.Tasks {
		t.Clear()
	}
}
