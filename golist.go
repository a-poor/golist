package golist

import "errors"

const IndentSize = 2

type TaskStatus int

const (
	TaskNotStarted TaskStatus = iota
	TaskInProgress
	TaskCompleted
	TaskFailed
	TaskSkipped
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

const (
	defaultTaskNotStarted = "➜"
	defaultTaskInProgress = "\\|/–"
	defaultTaskCompleted  = "✓"
	defaultTaskFailed     = "✗"
	defaultTaskSkipped    = "↓"
)

// ErrNilAction is returned when no action is
// set for a task.
var ErrNilAction = errors.New("nil action")
