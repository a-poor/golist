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

// Indicator is an interface for keeping track of
// status indicator state.
type Indicator interface {
	Get() string // Get the current indicator character
	Next()       // Move to the next indicator
}

// StaticIndicator implements the Indicator interface
// and returns a single (optionally colorized) indicator
// character.
//
// Note: The `Next` method is a no-op.
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

// CycleIndicator implements the Indicator interface and
// cycles through returning (optionally colorized)
// characters from a slice.
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
func (si *StatusIndicators) Get(s TaskStatus) string {
	backup := "–"
	i, ok := (*si)[s]
	if !ok {
		return backup
	}
	return i.Get()
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
//   – TaskNotStarted: "➜" (default terminal color)
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
