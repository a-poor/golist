package golist

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// IndentSize is the number of spaces to indent each line
// per task-depth level.
const IndentSize = 2

// Default List print delay
var DefaultListDelay = time.Millisecond * 100

var (
	// ErrListNotRunning is returned when certain functions (like `Println`)
	// are called that require the `List` to be running in order to work.
	ErrListNotRunning = errors.New("list not running")

	// ErrNoWriter is returned when the list is created without an `io.Writer`
	ErrNoWriter = errors.New("no writer specified")

	// ErrNilAction is returned when no action is set for a task
	ErrNilAction = errors.New("nil action")
)

// TaskStatus represents the current status of a task
type TaskStatus int

const (
	TaskNotStarted TaskStatus = iota // TaskNotStarted is the status for a task that hasn't started running yet
	TaskInProgress                   // TaskInProgress is the status for a task that is currently running
	TaskCompleted                    // TaskCompleted is the status for a task that has completed successfully
	TaskFailed                       // TaskFailed is the status for a task that returned a non-`nil` error
	TaskSkipped                      // TaskSkipped is the status for a task that was skipped (either manually or from a previous task's error)
)

// Format a TaskStatus as a string
func (s TaskStatus) String() string {
	switch s {
	case TaskNotStarted:
		return "Not Started"
	case TaskInProgress:
		return "In Progress"
	case TaskCompleted:
		return "Completed"
	case TaskFailed:
		return "Failed"
	case TaskSkipped:
		return "Skipped"
	default:
		return "Unknown"
	}
}

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
	printQ  chan string        // A channel for printing to the terminal while displaying the list
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
		return ErrNoWriter
	}

	// Check if it's already displaying
	if l.running {
		return nil
	}

	// Create a cancelable context
	ctx, cancel := context.WithCancel(context.Background())
	l.cancel = cancel

	// Create the channel for printing
	l.printQ = make(chan string)

	// Start the display loop
	go func() {
		// ts := l.getTaskStates()
		// l.print(ts)
		for {
			select {
			case <-ctx.Done():
				return
			case s := <-l.printQ:
				fmt.Fprintln(l.Writer, s)
			default:
				ts := l.getTaskStates()
				l.print(ts)
				l.StatusIniicator.Next()
				time.Sleep(time.Millisecond * 100)
				l.clear(ts)
			}
		}
	}()

	// Set the running flag
	l.running = true

	return nil
}

// Run starts running the tasks in the `List`
// and if `FailOnError` is set to true, returns
// an error if any of the tasks fail.
func (l *List) Run() error {
	l.Start()
	defer l.Stop()

	// Create a "base context"
	rootTaskCtx := taskContext{
		setMessage: func(m string) {},
		println: func(a ...interface{}) error {
			return l.Println(a...)
		},
		printfln: func(f string, a ...interface{}) error {
			return l.Printfln(f, a...)
		},
	}

	for _, t := range l.Tasks {
		if err := t.Run(&rootTaskCtx); err != nil && l.FailOnError {
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
	if l.cancel != nil {
		l.cancel()
	}
	if l.printQ != nil {
		close(l.printQ)
	}

	// Clear and print one final time (NOTE: should this be an option?)
	ts := l.getTaskStates()
	l.clear(ts)
	l.print(ts)

	l.running = false
	l.cancel = nil
	l.printQ = nil
}

// RunAndWait starts to display the task list statuses,
// runs the tasks, and waits for the tasks to complete
// before returning.
//
// RunAndWait is a convenience function that combines
// `Start`, `Run`, and `Stop`.
func (l *List) RunAndWait() error {
	err := l.Start()
	if err != nil {
		return err
	}
	err = l.Run()
	if err != nil {
		return err
	}
	l.Stop()
	return nil
}

// getTaskStates returns a slice of TaskStates
// for all child tasks
func (l *List) getTaskStates() []*TaskState {
	var messages []*TaskState
	for _, t := range l.Tasks {
		msgs := t.GetTaskStates()
		messages = append(messages, msgs...)
	}
	return messages
}

// truncateMessage will truncate the message to if it's too long
// based on the size parameter.
//
// If the message is truncated, all trailing spaces will be removed
// and an ellipsis ("…") is added to the end. An extra character
// will be removed to fit the elipsis, if necessary. If the size
//  is 0, an ellipsis character is still returned.
func (l *List) truncateMessage(m string, size int) string {
	rm := []rune(m)
	if len(rm) <= size { // No truncation needed
		return m
	}
	if size <= 1 { // Truncate everything
		return "…"
	}

	// Remove an extra character to fit the ellipsis
	tsize := size - 1

	return strings.TrimSuffix(string(rm[0:tsize]), " ") + "…"
}

// formatMessage formats a message row for displaying.
// The format used is: [depth] [status] [message]
// and it's length is (optionally) limited by the
// MaxLineLength parameter.
func (l *List) formatMessage(m *TaskState) string {
	n := m.Depth * IndentSize
	d := strings.Repeat(" ", n)
	i := l.StatusIniicator.Get(m.Status)

	// If no no truncate text, just return the formatted
	// status message
	if l.MaxLineLength == 0 {
		return fmt.Sprintf("%s%s %s", d, i, m.Message)
	}

	// Otherwise, truncate the result
	size := l.MaxLineLength - (n + 1)
	return fmt.Sprintf("%s%s %s", d, i, l.truncateMessage(m.Message, size))
}

// print prints the current task states
func (l *List) print(states []*TaskState) {
	for _, m := range states {
		fmt.Fprintln(l.Writer, l.formatMessage(m))
	}
}

// print prints the current task states
func (l *List) clear(states []*TaskState) {
	n := len(states)
	s := "\033[1A" // Move up a line
	s += "\033[K"  // Clear the line
	s += "\r"      // Move back to the beginning of the line

	for i := 0; i < n; i++ {
		fmt.Fprint(l.Writer, s)
	}
}

// Println prints information to the List's Writer (which is
// most likely stdout), like `fmt.Println`.
//
// Internally, it passes information to a channel that will
// be read by the display goroutine and printed safely in
// between updates to the task-list.
//
// Note: If Println is called while the list is not running,
// it will return the error ErrListNotRunning.
func (l *List) Println(a ...interface{}) error {
	if l.printQ == nil {
		return ErrNoWriter
	}
	s := fmt.Sprint(a...)
	l.printQ <- s
	return nil
}

// Printfln prints a formatted string to the list's writer
// (which is usually stdout) before reprinting the list.
//
// Printfln is like Printf but adds a newline character at
// the end of the string. It requires that there be a newline
// character so that the list can be reprinted properly.
//
// Note: If Printfln is called while the list is not running,
// it will return the error ErrListNotRunning.
func (l *List) Printfln(f string, d ...interface{}) error {
	if l.printQ == nil {
		return ErrNoWriter
	}
	s := fmt.Sprintf(f, d...)
	l.printQ <- s
	return nil
}
