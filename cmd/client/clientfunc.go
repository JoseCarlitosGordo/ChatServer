package main

import (
	"bufio"
	extras "chatserver/structs"
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

var msgScanner = bufio.NewScanner(os.Stdin)

func sendMessage() string {

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

func CommenceAuthenticationProcess(connection net.Conn) (extras.Connection, error) {
	//TODO:: Prompt for login, sign up or guest. Then
	fmt.Print("1. To Log in, press '1' \n 2. To Sign Up, press '2' \n To register as a guest, press any other key ")
	msgScanner.Scan()
	input := msgScanner.Text()
	if input == "1" {
		//return LoginProcess()
	}
	if input == "2" {
		//return SignUpProcess()
	}
	fmt.Println("Logging in as Guest")
	return extras.Connection{ConnectionObj: connection, Account: extras.UserAccount{}}, nil

}

func LoginProcess(connection net.Conn) (extras.Connection, error) {
	//TODO: pass in connection, check for user account email and password on server and return a valid Connection Object\

	// Create an encoder and target our buffer
	enc := gob.NewEncoder(connection)

	fmt.Println("Logging in")
	fmt.Print("Your username: ")
	msgScanner.Scan()
	userName := msgScanner.Text()

	fmt.Print("Your Password: ")
	msgScanner.Scan()
	password := msgScanner.Text()

	loginAttempt := extras.UserAccount{UserName: userName, Password: password}

	//send login attempt over network
	err := enc.Encode(extras.Packet[extras.Connection]{Type: "Login Attempt", Content: extras.Connection{ConnectionObj: connection, Account: loginAttempt}})
	if err != nil {
		return extras.Connection{}, err
	}
	return extras.Connection{}, nil
	//TODO: read from buffer to check if username and password match what is found in the db. If it is found, return True.
	//Otherwise, continue login process

}

func SignUpProcess() {
	//TODO: pass in connection, prompt user for user name and password, hash the password, salt it, return a valid connection object.
	//This connection obj will contain user account and the tcp connection to the server, ensuring both client and server knows who it belongs to
}
