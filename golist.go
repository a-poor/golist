package main

import (
	"fmt"
	"time"
)

const (
	EmojiCheck    = "✔️"
	EmojiX        = "✖"
	EmojiXAlt     = "✗"
	ArrowRight    = "➜"
	ArrowRightAlt = "❱"
)

type TaskStatus int

const (
	TaskNotStarted TaskStatus = iota
	TaskInProgress
	TaskCompleted
	TaskFailed
)

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

func main() {
	// fmt.Print("Starting")
	// time.Sleep(time.Second * 2)
	// ClearScreen()
	// fmt.Println("Done")

	i := 0
	spinnerChars := []rune("⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏")
	words := []string{"Hello", "World", "!"}

	for {
		// fmt.Printf("\r%c %s", spinnerChars[i], words[i%len(words)])
		fmt.Printf("\033[1K\r%c %s", spinnerChars[i], words[i%len(words)])

		time.Sleep(time.Millisecond * 100)
		// time.Sleep(time.Second * 1)

		i = (i + 1) % len(spinnerChars)
	}

}
