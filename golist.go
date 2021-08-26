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
		return fmtBlack("➜")
	case TaskInProgress:
		return fmtYellow(s.spinner.GetAndIncrement())
	case TaskCompleted:
		return fmtGreen("✔️")
	case TaskFailed:
		return fmtRed("✗")
	default:
		return "-"
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
		for {
			s.PrintStatus()
			time.Sleep(time.Millisecond * 100)
			if stat := s.Get(); stat == TaskCompleted || stat == TaskFailed {
				s.PrintStatus()
				fmt.Println()
				break
			}
		}
		time.Sleep(time.Second)
	}()
	time.Sleep(time.Second)
	s.Set(TaskInProgress)
	time.Sleep(time.Second)
	if rand.Float64() > 0.5 {
		s.Set(TaskCompleted)
	} else {
		s.Set(TaskFailed)
	}
	time.Sleep(time.Millisecond * 100)
}
