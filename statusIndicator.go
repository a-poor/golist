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
	defaultTaskNotStarted = "➜"
	defaultTaskInProgress = "🕐🕑🕒🕓🕔🕕🕖🕗🕘🕙🕚🕛"
	defaultTaskCompleted  = "✓"
	defaultTaskFailed     = "✗"
	defaultTaskSkipped    = "↓"
)

func toBlack(s string) string {
	return "\033[30;1m" + s + "\033[0m"
}

func toYellow(s string) string {
	return "\033[33;1m" + s + "\033[0m"
}

func toGreen(s string) string {
	return "\033[32;1m" + s + "\033[0m"
}

func toRed(s string) string {
	return "\033[31;1m" + s + "\033[0m"
}

type statusIndicatorConfig struct {
	NotStarted string
	InProgress string
	Completed  string
	Failed     string
	Skipped    string
}

func createStatusIndicator(config statusIndicatorConfig) func(TaskStatus) string {
	// Create maps to store state
	indexes := make(map[TaskStatus]int)
	characters := make(map[TaskStatus][]rune)

	// Get values (or use defaults)
	if config.NotStarted == "" {
		characters[TaskNotStarted] = []rune(defaultTaskNotStarted)
	} else {
		characters[TaskNotStarted] = []rune(config.NotStarted)
	}
	if config.InProgress == "" {
		characters[TaskInProgress] = []rune(defaultTaskInProgress)
	} else {
		characters[TaskInProgress] = []rune(config.InProgress)
	}
	if config.Completed == "" {
		characters[TaskCompleted] = []rune(defaultTaskCompleted)
	} else {
		characters[TaskCompleted] = []rune(config.Completed)
	}
	if config.Failed == "" {
		characters[TaskFailed] = []rune(defaultTaskFailed)
	} else {
		characters[TaskFailed] = []rune(config.Failed)
	}
	if config.Skipped == "" {
		characters[TaskSkipped] = []rune(defaultTaskSkipped)
	} else {
		characters[TaskSkipped] = []rune(config.Skipped)
	}

	// Create & return closure
	return func(status TaskStatus) string {
		// Get index
		i := indexes[status]

		// Get character
		cs, ok := characters[status]
		if !ok {
			return " "
		}

		var colorFn func(string) string
		switch status {
		case TaskNotStarted, TaskSkipped:
			colorFn = toBlack
		case TaskInProgress:
			colorFn = toYellow
		case TaskCompleted:
			colorFn = toGreen
		case TaskFailed:
			colorFn = toRed
		default:
			colorFn = func(s string) string {
				return s
			}
		}

		// Increment index
		indexes[status] = (i + 1) % len(cs)

		// Return character
		return colorFn(string(cs[i]))
	}
}
