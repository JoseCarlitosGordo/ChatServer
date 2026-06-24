package main

import (
	"bufio"
	"fmt"
	"os"
)

func sendMessage() string {
	msgScanner := bufio.NewScanner(os.Stdin)
	msgScanner.Scan()
	message := msgScanner.Text()
	// 1. Move cursor up 1 line
	fmt.Print("\033[1A")
	// ANSI Escape Sequence to clear the entire current line
	fmt.Print("\033[2K")

	// Carriage return to move cursor to the beginning of the line
	fmt.Print("\r")

	return message
}
