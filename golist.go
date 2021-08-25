package main

import (
	"fmt"
	"time"
)

func SetCursor(line, column int) {
	fmt.Printf("\033[%d;%dH", line, column)
}

func MoveCursorUp(n int) {
	fmt.Printf("\033[%dA", n)
}

func MoveCursorDown(n int) {
	fmt.Printf("\033[%dB", n)
}

func MoveCursorForward(n int) {
	fmt.Printf("\033[%dC", n)
}

func MoveCursorBackward(n int) {
	fmt.Printf("\033[%dD", n)
}

func ClearScreen() {
	fmt.Print("\033[2J")
}

func EraseToEndOfLine() {
	fmt.Print("\033[K")
}

func SaveCursorPosition() {
	fmt.Print("\033[s")
}

func RestoreCursorPosition() {
	fmt.Print("\033[u")
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
		fmt.Printf("\r%c %s", spinnerChars[i], words[i%len(words)])
		// fmt.Printf("\033[1K\r%c %s", spinnerChars[i], words[i%len(words)])

		// time.Sleep(time.Millisecond * 100)
		time.Sleep(time.Second * 1)

		i = (i + 1) % len(spinnerChars)
	}

}
