package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	EmojiCheck    = "✔️"
	EmojiX        = "✖"
	EmojiXAlt     = "✗"
	ArrowRight    = "➜"
	ArrowRightAlt = "❱"
)

const SpinnerChars = "⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏"

type TaskStatus int

const (
	TaskNotStarted TaskStatus = iota
	TaskInProgress
	TaskCompleted
	TaskFailed
)

var StatusMsg = map[TaskStatus]string{
	TaskNotStarted: "Not Started",
	TaskInProgress: "In Progress",
	TaskCompleted:  "Completed",
	TaskFailed:     "Failed",
}

type List struct {
	Tasks []Task
}

type Task struct {
	Name     string
	Subtasks []Task
}

type Spinner struct {
	Chars        []rune
	currentIndex int
}

func (s *Spinner) Get() string {
	return string(s.Chars[s.currentIndex])
}

func (s *Spinner) Increment() {
	s.currentIndex = (s.currentIndex + 1) % len(s.Chars)
}

func (s *Spinner) GetAndIncrement() string {
	c := s.Get()
	s.Increment()
	return c
}

type status struct {
	sync.RWMutex
	status  TaskStatus
	spinner Spinner
}

func (s *status) Get() TaskStatus {
	s.RLock()
	defer s.RUnlock()
	return s.status
}

func (s *status) Set(status TaskStatus) {
	s.Lock()
	defer s.Unlock()
	s.status = status
}

func (s *status) GetStatusChar(status TaskStatus) string {
	switch status {
	case TaskNotStarted:
		return "\033[30;1m➜\033[0m"
	case TaskInProgress:
		return "\033[33;1m" + s.spinner.GetAndIncrement() + "\033[0m"
	case TaskCompleted:
		return "\033[32;1m✔️\033[0m"
	case TaskFailed:
		return "\033[31;1m✗\033[0m"
	default:
		return "?"
	}
}

func (s *status) GetMessage() string {
	stat := s.Get()
	return fmt.Sprintf("\033[1K\r%s %s", s.GetStatusChar(stat), StatusMsg[stat])
}

func (s *status) PrintStatus() {
	fmt.Print(s.GetMessage())
}

func main() {
	rand.Seed(time.Now().UnixNano())

	s := status{
		spinner: Spinner{Chars: []rune(SpinnerChars)},
	}

	go func() {
		time.Sleep(time.Second)
		s.Set(TaskInProgress)
		time.Sleep(time.Second)
		if rand.Float64() > 0.5 {
			s.Set(TaskCompleted)
		} else {
			s.Set(TaskFailed)
		}
	}()

	for {
		s.PrintStatus()
		time.Sleep(time.Millisecond * 100)
		if stat := s.Get(); stat == TaskCompleted || stat == TaskFailed {
			s.PrintStatus()
			fmt.Println()
			break
		}
	}
	// time.Sleep(time.Second)
}
