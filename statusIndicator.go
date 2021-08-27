package main

type TaskStatus int

const (
	TaskNotStarted TaskStatus = iota
	TaskInProgress
	TaskCompleted
	TaskFailed
	TaskSkipped
)

const (
	DefaultTaskNotStarted = "➜"
	DefaultTaskInProgress = "🕐🕑🕒🕓🕔🕕🕖🕗🕘🕙🕚🕛"
	DefaultTaskCompleted  = "✓"
	DefaultTaskFailed     = "✗"
	DefaultTaskSkipped    = "↓"
)

type StatusIndicatorConfig struct {
	NotStarted string
	InProgress string
	Completed  string
	Failed     string
	Skipped    string
}

func CreateStatusIndicator(config *StatusIndicatorConfig) func(TaskStatus) rune {
	// Create maps to store state
	indexes := make(map[TaskStatus]int)
	characters := make(map[TaskStatus][]rune)

	// Get values (or use defaults)
	if config.NotStarted == "" {
		characters[TaskNotStarted] = []rune(DefaultTaskNotStarted)
	} else {
		characters[TaskNotStarted] = []rune(config.NotStarted)
	}
	if config.InProgress == "" {
		characters[TaskInProgress] = []rune(DefaultTaskInProgress)
	} else {
		characters[TaskInProgress] = []rune(config.InProgress)
	}
	if config.Completed == "" {
		characters[TaskCompleted] = []rune(DefaultTaskCompleted)
	} else {
		characters[TaskCompleted] = []rune(config.Completed)
	}
	if config.Failed == "" {
		characters[TaskFailed] = []rune(DefaultTaskFailed)
	} else {
		characters[TaskFailed] = []rune(config.Failed)
	}
	if config.Skipped == "" {
		characters[TaskSkipped] = []rune(DefaultTaskSkipped)
	} else {
		characters[TaskSkipped] = []rune(config.Skipped)
	}

	// Create & return closure
	return func(status TaskStatus) rune {
		// Get index
		i := indexes[status]

		// Get character
		cs, ok := characters[status]
		if !ok {
			return ' '
		}

		// Increment index
		indexes[status] = (i + 1) % len(cs)

		// Return character
		return cs[i]
	}
}
