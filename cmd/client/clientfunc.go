package main

import (
	"bufio"
	"fmt"
	"os"
)

func sendMessage() string {
	msgScanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("\n(Type 'exit()' to close this app) \nYour Message: ")
	msgScanner.Scan()
	message := msgScanner.Text()
	return message
}
