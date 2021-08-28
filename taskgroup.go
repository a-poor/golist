package main

import (
	"fmt"
	"strings"
)

type TaskGroup struct {
	Message     string       // The message to be displayed
	Tasks       []TaskRunner // A list of tasks to run
	Skip        func() bool  // Is run before the task starts. If returns true, the task isn't run
	FailOnError bool         // If true, the task group stops on the first error

	statusIndicator func(TaskStatus) string
	err             error      // The error that occurred during the last task
	status          TaskStatus // The status of the task
}

func NewTaskGroup() *TaskGroup {
	return &TaskGroup{}
}

func (tg *TaskGroup) AddTask(t *TaskRunner) {
	tg.Tasks = append(tg.Tasks, *t)
}

func (tg *TaskGroup) Run() error {
	if tg.statusIndicator == nil {
		tg.initStatusIndicator()
	}

	// Check if the task should be skipped
	if tg.Skip != nil && tg.Skip() {
		tg.SetStatus(TaskSkipped)
		return nil
	}

	// Run the tasks
	tg.SetStatus(TaskInProgress)
	var err error
	for _, t := range tg.Tasks {
		err = t.Run()
		if err != nil && tg.FailOnError {
			break
		}
	}

	// Set the final status
	if err != nil {
		tg.SetStatus(TaskFailed)
	} else {
		tg.SetStatus(TaskCompleted)
	}
	return err
}

func (tg *TaskGroup) Print(indent int) {
	pad := strings.Repeat(" ", indent)
	stat := tg.statusIndicator(tg.status)
	fmt.Printf("%s%s %s\n", pad, stat, tg.Message)
	for _, t := range tg.Tasks {
		t.Print(indent + IndentSize)
	}
}

func (tg *TaskGroup) Clear() {
	// Clear the subtasks in reverse order
	for i := len(tg.Tasks) - 1; i >= 0; i-- {
		// Clear the subtasks
		tg.Tasks[i].Clear()
	}

	// Then clear the task's line
	fmt.Print("\033[1A") // Move up a line
	fmt.Print("\033[K")  // Clear the line
	fmt.Print("\r")      // Move back to the beginning of the line
}

func (tg *TaskGroup) SetMessage(m string) {
	tg.Message = m
}

func (tg *TaskGroup) GetError() error {
	return tg.err
}

func (tg *TaskGroup) SetError(err error) {
	tg.err = err
}

func (tg *TaskGroup) SetStatus(s TaskStatus) {
	tg.status = s
}

func (tg *TaskGroup) GetSize() int {
	return len(tg.Tasks) + 1
}

func (tg *TaskGroup) initStatusIndicator() {
	tg.statusIndicator = createStatusIndicator(&statusIndicatorConfig{})
}
