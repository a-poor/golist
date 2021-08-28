package main

import (
	"fmt"
	"strings"
)

type Task struct {
	Message string       // Message to display to user
	Action  func() error // The task function to be run
	Skip    func() bool  // Is run before the task starts. If returns true, the task isn't run

	statusIndicator func(TaskStatus) string
	status          TaskStatus // The status of the task
	err             error      // The error returned by the task function

}

func (t *Task) Run() error {
	if t.statusIndicator == nil {
		t.initStatusIndicator()
	}
	return nil
}

func (t *Task) Print(indent int) {
	pad := strings.Repeat(" ", indent)
	stat := t.statusIndicator(t.status)
	fmt.Printf("%s%s %s\n", pad, stat, t.Message)
}

func (t *Task) Clear() {
	fmt.Print("\033[1A") // Move up a line
	fmt.Print("\033[K")  // Clear the line
	fmt.Print("\r")      // Move back to the beginning of the line
}

func (t *Task) SetMessage(m string) {
	t.Message = m
}

func (t *Task) SetError(err error) {
	t.err = err
}

func (t *Task) SetStatus(s TaskStatus) {
	t.status = s
}

func (t *Task) GetDepth() int {
	return 1
}

func (t Task) GetSize() int {
	return 1
}

func (t *Task) initStatusIndicator() {
	t.statusIndicator = createStatusIndicator(&statusIndicatorConfig{})
}
