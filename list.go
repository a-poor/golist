package main

import (
	"context"
	"fmt"
	"time"
)

const (
	EmojiCheck    = "✔️"
	EmojiX        = "✖"
	EmojiXAlt     = "✗"
	ArrowRight    = "➜"
	ArrowRightAlt = "❱"
	ArrowDown     = "↓"
)

const (
	SpinnerChars1 = "⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏"
	SpinnerChars2 = "🕐🕑🕒🕓🕔🕕🕖🕗🕘🕙🕚🕛"
	SpinnerChars3 = "|/-\\"
)

type TaskStatus int

const (
	TaskNotStarted TaskStatus = iota
	TaskInProgress
	TaskCompleted
	TaskFailed
	TaskSkipped
)

type List struct {
	Tasks   []*Task
	running bool
	cancel  context.CancelFunc
}

func NewList() *List {
	return &List{Tasks: make([]*Task, 0)}
}

func (l *List) AddTask(t *Task) {
	l.Tasks = append(l.Tasks, t)
}

func (l *List) AddNewTask(title string, action func() error) {
	t := NewTask(title, action)
	l.AddTask(t)
}

func (l *List) Run() error {
	l.Start()
	defer l.Stop()
	for _, t := range l.Tasks {
		if err := t.Run(); err != nil && !t.ContinueOnFail {
			return err
		}
	}
	return nil
}

func (l *List) print() {
	for _, t := range l.Tasks {
		fmt.Printf("\r%s\n", t.String())
	}
}

func (l *List) clear() {
	nLines := len(l.Tasks)
	for i := 0; i < nLines; i++ {
		fmt.Print("\033[1A")
		fmt.Print("\033[K")
	}
}

// Start displaying the list from a goroutine
func (l *List) Start() {
	if l.running {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	l.cancel = cancel
	go func() {
		l.print()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				l.clear()
				l.print()
				time.Sleep(time.Millisecond * 100)
			}
		}
	}()
	l.running = true
}

func (l *List) Stop() {
	// Confirm that the list is running
	if !l.running {
		return
	}

	// Send the cancel signal
	l.cancel()

	// Print the final results
	l.clear()
	l.print()

	// Set the running flag and clear cancel function
	l.running = false
	l.cancel = nil
}

type spinner struct {
	chars []rune
	index int
}

func NewSpinner() *spinner {
	return &spinner{chars: []rune(SpinnerChars3)}
}

func (s *spinner) Next() string {
	c := s.chars[s.index]
	s.index = (s.index + 1) % len(s.chars)
	return string(c)
}

type Task struct {
	spin           *spinner
	Status         TaskStatus
	Title          string
	Skip           func() bool  // Runs before the task. Should it run?
	Action         func() error // The actual function
	ContinueOnFail bool
	Err            error
}

func NewTask(title string, action func() error) *Task {
	return &Task{
		spin:   NewSpinner(),
		Title:  title,
		Action: action,
	}
}

func (t *Task) Run() error {
	// Check if the task should be skipped
	if t.Skip != nil && t.Skip() {
		t.Status = TaskSkipped
		return nil
	}

	// Set the status to be "in progress"
	t.Status = TaskInProgress

	// Run the task
	err := t.Action()

	// Store the error
	t.Err = err
	if err != nil {
		t.Status = TaskFailed
	} else {
		t.Status = TaskCompleted
	}
	return err
}

func (t *Task) getMarker() string {
	switch t.Status {
	case TaskNotStarted:
		return toBlack(ArrowRight)
		// return toBlack(t.spin.Next())
	case TaskInProgress:
		return toYellow(t.spin.Next())
	case TaskCompleted:
		return toGreen(EmojiCheck)
	case TaskFailed:
		return toRed(EmojiX)
	case TaskSkipped:
		return toBlack(ArrowDown)
	default:
		return "?"
	}
}

func (t *Task) String() string {
	return fmt.Sprintf("%s %s", t.getMarker(), t.Title)
}
