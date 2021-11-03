package golist

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
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

// List is the top-level list object that
// represents a group of tasks to be run.
//
// Generally, you'll want to use the `NewList`
// to create a new list with some sensible defaults.
//
// If `Writer` is not set, it will be set to `os.Stdout`.
// If `StatusIndicator` is not set it will default
// to using the "–" character for all list items.
//
// Typically, you'll at least want to set `Writer`,
// `Delay`, and `StatusIndicator`.
type List struct {
	Writer          io.Writer        // Writer to use for printing output
	Delay           time.Duration    // Delay between prints
	StatusIndicator StatusIndicators // Map of statuses to status indicators
	Tasks           []TaskRunner     // List of tasks to run
	FailOnError     bool             // If true, the task execution stops on the first error. Note: this will be ignored if Concurrent is true.
	MaxLineLength   int              // Maximum line length for printing (0 = no limit)
	ClearOnComplete bool             // If true, the list will clear the list after it finishes running
	Concurrent      bool             // Should the tasks be run concurrently? Note: If true, ignores the FailOnError flag

	printDone chan bool          // Is the printing loop done
	running   bool               // Is the list running?
	cancel    context.CancelFunc // A context cancel function for stopping the list run
	printQ    chan string        // A channel for printing to the terminal while displaying the list

	writeLock sync.Mutex
}

// NewList creates a new task list with some sensible defaults.
// It writes to stdout and and has a delay of 100ms between prints.
func NewList() *List {
	return &List{
		Writer:          os.Stdout,
		Delay:           DefaultListDelay,
		StatusIndicator: CreateDefaultStatusIndicator(),
	}
}

// NewList creates a new task list that writes to the
// provided io.Writer. Mostly used for testing.
func NewListWithWriter(w io.Writer) *List {
	return &List{
		Writer:          w,
		Delay:           DefaultListDelay,
		StatusIndicator: CreateDefaultStatusIndicator(),
	}
}

// AddTask adds a TaskRunner to the top-level List
// and returns a pointer to itself.
func (l *List) AddTask(t TaskRunner) *List {
	if l.Tasks == nil {
		l.Tasks = make([]TaskRunner, 0)
	}
	l.Tasks = append(l.Tasks, t)
	return l
}

// Start begins displaying the list statuses
// from a background goroutine.
//
// Note: If the list is created without a writer,
// it will be set to `os.Stdout`.
func (l *List) Start() {
	l.writeLock.Lock()
	if l.Writer == nil {
		l.Writer = os.Stdout
	}
	l.writeLock.Unlock()

	// Check if it's already displaying
	if l.running {
		return
	}

	// Create a cancelable context
	ctx, cancel := context.WithCancel(context.Background())
	l.cancel = cancel

	// Create the channel for printing
	l.printQ = make(chan string)

	// Create a channel to tel the Stop function when the
	// print loop has completed
	l.printDone = make(chan bool)
	donePrinting := func() {
		l.printDone <- true
	}

	// Set the running flag
	l.running = true

	// Start the display loop
	go func() {
		defer donePrinting() // Tell the Stop function that we're done printing
		ts := l.getTaskStates()
		l.print(ts)
		nLinesToClear := len(ts)
		for {
			select {
			case <-ctx.Done(): // Check if the print loop should stop

				// Perform a final clear and an optional print
				// depending on `ClearOnComplete`
				l.clear(nLinesToClear)
				if l.ClearOnComplete {
					return
				}

				ts := l.getTaskStates()
				l.print(ts)
				return

			case s := <-l.printQ: // Check if there's a message to print
				l.clear(nLinesToClear)
				nLinesToClear = 0
				l.writeLock.Lock()
				fmt.Fprintln(l.Writer, s)
				l.writeLock.Unlock()

			default:
				l.clear(nLinesToClear)

				ts := l.getTaskStates()
				l.print(ts)
				l.StatusIndicator.Next()

				nLinesToClear = len(ts)

				time.Sleep(l.Delay)
			}
		}
	}()
}

// createRootContext creates a base TaskContext to be passed
// down to the subtasks to create sub-contexts.
//
// Note: The SetMessage function is a no-op, since the
// top-level list doesn't have a message to set.
func (l *List) createRootContext() TaskContext {
	return &taskContext{
		setMessage: func(m string) {},
		println: func(a ...interface{}) error {
			return l.Println(a...)
		},
		printfln: func(f string, a ...interface{}) error {
			return l.Printfln(f, a...)
		},
	}
}

// runSync runs the TaskRunners in this TaskGroup synchronously
func (l *List) runSync(c TaskContext) error {
	var skipRemaining bool
	for _, t := range l.Tasks {
		if skipRemaining {
			t.SetStatus(TaskSkipped)
			continue
		}
		err := t.Run(c)
		if err != nil && l.FailOnError {
			skipRemaining = true
		}
	}
	return l.GetError()
}

// runAsync runs the TaskRunners in this TaskGroup concurrently
// and blocks until they are all done.
func (l *List) runAsync(c TaskContext) error {
	var wg sync.WaitGroup
	for _, t := range l.Tasks {
		wg.Add(1)
		go func(t TaskRunner, c TaskContext) {
			defer wg.Done()
			t.Run(c)
		}(t, c)
	}
	wg.Wait()
	time.Sleep(l.Delay)
	return l.GetError()
}

// Run starts running the tasks in the `List`
// and if `FailOnError` is set to true, returns
// an error if any of the tasks fail.
func (l *List) Run() error {
	// Starts the list if it hasn't already started
	l.Start()

	// Create a "base context" to be passed down
	// for subtasks to create TaskContexts
	rootTaskCtx := l.createRootContext()

	// Check if running concurrently...
	var err error
	if l.Concurrent {
		err = l.runAsync(rootTaskCtx)
	} else {
		// Otherwise, run synchronously
		err = l.runSync(rootTaskCtx)
	}

	// Return the error
	return err
}

// Stop stops displaying the task list statuses and
// cancels the background goroutine.
//
// Stop also clears and prints one final time
// before finishing.
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

	// Wait for the print loop to finish
	<-l.printDone

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
	l.Start()
	err := l.Run()
	l.Stop()
	return err
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
	i := l.StatusIndicator.Get(m.Status)

	// If no no truncate text, just return the formatted
	// status message
	if l.MaxLineLength == 0 {
		return fmt.Sprintf("%s%s %s", d, i, m.Message)
	}

	// Otherwise, truncate the result
	size := l.MaxLineLength - (n + 1)
	return fmt.Sprintf("%s%s %s", d, i, l.truncateMessage(m.Message, size))
}

// fmtPrint returns the formatted list of messages
// and statuses, using the supplied TaskStates
func (l *List) fmtPrint(ts []*TaskState) string {
	s := make([]string, 0)
	for _, t := range ts {
		s = append(s, l.formatMessage(t))
	}
	return strings.Join(s, "\n")
}

// print prints the current task states
func (l *List) print(states []*TaskState) {
	s := l.fmtPrint(states)
	fmt.Fprintln(l.Writer, s)
}

// fmtClear returns a string of ANSI escape characters
// to clear the `n` lines previously printed.
func (l *List) fmtClear(n int) string {
	s := "\033[1A" // Move up a line
	s += "\033[K"  // Clear the line
	s += "\r"      // Move back to the beginning of the line
	return strings.Repeat(s, n)
}

// clear clears the previous task states using
// ANSII escape characters
func (l *List) clear(n int) {
	s := l.fmtClear(n)
	fmt.Fprint(l.Writer, s)
}

// clearThenPrint is a shorcut that clears the previous
// task states and prints the current task states.
//
// It is equivalent to calling `clear` and `print`
// except that it only uses one call to `fmt.Fprintln`.
func (l *List) clearThenPrint(states []*TaskState) {
	c := l.fmtClear(len(states))
	s := l.fmtPrint(states)
	fmt.Fprintln(l.Writer, c+s)
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

// GetError returns the errors from the child tasks
func (l *List) GetError() error {
	var err *multierror.Error
	for _, t := range l.Tasks {
		err = multierror.Append(err, t.GetError())
	}
	return err.ErrorOrNil()
}
