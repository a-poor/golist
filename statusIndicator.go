package golist

// ToBlack wraps a string in escape characters
// to format it as black text.
func ToBlack(s string) string {
	return "\033[30;1m" + s + "\033[0m"
}

// ToYellow wraps a string in escape characters
// to format it as yellow text.
func ToYellow(s string) string {
	return "\033[33;1m" + s + "\033[0m"
}

// ToGreen wraps a string in escape characters
// to format it as green text.
func ToGreen(s string) string {
	return "\033[32;1m" + s + "\033[0m"
}

// ToRed wraps a string in escape characters
// to format it as red text.
func ToRed(s string) string {
	return "\033[31;1m" + s + "\033[0m"
}

type Indicator interface {
	Get() string // Get the current indicator character
	Next()       // Move to the next indicator
}

type StaticIndicator struct {
	Indicator rune                // Character to return
	Colorizer func(string) string // Optional function to colorize the indicator
}

// Get returns the current status indicator.
// If Colorizer is set, calls it on the indicator character.
func (si *StaticIndicator) Get() string {
	s := string(si.Indicator)
	if si.Colorizer == nil {
		return s
	}
	return si.Colorizer(s)
}

// Next is a no-op for StaticIndicator
func (si *StaticIndicator) Next() {
	// Do nothing
}

type CycleIndicator struct {
	Indicators []rune              // Array of indicator characters
	index      int                 // Current position in the Indicators array
	Colorizer  func(string) string // Optional function to colorize the indicator
}

// Get returns the current status indicator.
// If Colorizer is set, calls it on the indicator character.
func (si *CycleIndicator) Get() string {
	s := string(si.Indicators[si.index])
	if si.Colorizer == nil {
		return s
	}
	return si.Colorizer(s)
}

// Next increments the current index in the Indicators array
// and wraps around if passed the end of the array.
func (si *CycleIndicator) Next() {
	si.index = (si.index + 1) % len(si.Indicators)
}

// StatusIndicators is a map from task statuses to an indicator.
type StatusIndicators map[TaskStatus]Indicator

// Get returns the current status indicator character
// for the corresponding TaskStatus in the StatusIndicators map.
//
// Note: If the TaskStatus is not found in the StatusIndicators map,
// the uncolorized string "–" is returned
func (si *StatusIndicators) Get() string {
	backup := "–"
	s, ok := (*si)[TaskNotStarted]
	if !ok {
		return backup
	}
	return s.Get()
}

// Next calls Next on all indicators
func (si *StatusIndicators) Next() {
	for _, i := range *si {
		i.Next()
	}
}

// CreateDefaultStatusIndicator creates a StatusIndicators map
// with default values for each status.
//
// The default values are:
//   – TaskNotStarted: "➜"
//   – TaskInProgress: "⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏" (yellow)
//   – TaskCompleted: "✓" (green)
//   – TaskFailed: "✗" (red)
//   – TaskSkipped: "↓" (black)
//
func CreateDefaultStatusIndicator() StatusIndicators {
	return StatusIndicators{
		TaskNotStarted: &StaticIndicator{
			Indicator: '➜',
		},
		TaskInProgress: &CycleIndicator{
			Indicators: []rune("⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏"),
			Colorizer:  ToYellow,
		},
		TaskCompleted: &StaticIndicator{
			Indicator: '✓',
			Colorizer: ToGreen,
		},
		TaskFailed: &StaticIndicator{
			Indicator: '✗',
			Colorizer: ToRed,
		},
		TaskSkipped: &StaticIndicator{
			Indicator: '↓',
			Colorizer: ToBlack,
		},
	}
}

//////// OLD FUNCTION ////////////

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
			colorFn = ToBlack
		case TaskInProgress:
			colorFn = ToYellow
		case TaskCompleted:
			colorFn = ToGreen
		case TaskFailed:
			colorFn = ToRed
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
