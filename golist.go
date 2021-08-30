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

const (
	defaultTaskNotStarted = "➜"
	defaultTaskInProgress = "\\|/–"
	defaultTaskCompleted  = "✓"
	defaultTaskFailed     = "✗"
	defaultTaskSkipped    = "↓"
)

var ErrNilAction = errors.New("nil action")

var StatusMsg = map[TaskStatus]string{
	TaskNotStarted: "Not Started",
	TaskInProgress: "In Progress",
	TaskCompleted:  "Completed",
	TaskFailed:     "Failed",
	TaskSkipped:    "Skipped",
}
