package main

import (
	"bufio"
	"fmt"
	"net"
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

func CommenceAuthenticationProcess(connection net.Conn) (bool, error) {
	//TODO:: Prompt for login, sign up or guest. Then
	msgScanner := bufio.NewScanner(os.Stdin)
	fmt.Print("1. To Log in, press '1' \n 2. To Sign Up, press '2' \n To register as a guest, press any other key ")
	msgScanner.Scan()
	input := msgScanner.Text()
	if input == "1" {
		//return LoginProcess()
	}
	if input == "2" {
		//return SignUpProcess()
	}
	return false, nil

}
