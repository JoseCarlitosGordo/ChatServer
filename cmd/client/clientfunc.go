package main

import (
	"bufio"
	extras "chatserver/structs"
	"encoding/gob"
	"fmt"
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

func CommenceAuthenticationProcess(sessionState *extras.SessionState) error {
	//TODO:: Prompt for login, sign up or guest. Then
	fmt.Print("1. To Log in, press '1' \n 2. To Sign Up, press '2' \n To register as a guest, press any other key ")
	msgScanner.Scan()
	input := msgScanner.Text()
	if input == "1" {
		return LoginProcess(sessionState)
	}
	if input == "2" {
		return SignUpProcess(sessionState)
	}
	fmt.Println("Logging in as Guest")
	return nil

}

func LoginProcess(sessionState *extras.SessionState) error {
	//TODO: pass in connection, check for user account email and password on server and return a valid Connection Object\

	// Create an encoder and target our buffer
	enc := gob.NewEncoder(sessionState.ConnectionWrapper.ConnectionObj)
	for {

		fmt.Println("Logging in")
		fmt.Print("Your username: ")
		msgScanner.Scan()
		userName := msgScanner.Text()

		fmt.Print("Your Password: ")
		msgScanner.Scan()
		password := msgScanner.Text()

		loginAttempt := extras.UserAccount{UserName: userName, Password: password}

		//send login attempt over network
		err := enc.Encode(extras.Packet[extras.Connection]{Type: "Login Attempt", Content: extras.Connection{ConnectionObj: sessionState.ConnectionWrapper.ConnectionObj, Account: loginAttempt}})
		if err != nil {
			return err
		}

		//TODO: refactor client so that it handles messages and other packets sent from the server
		serverResponse := <-sessionState.PacketListener
		response := serverResponse.Content.(extras.Packet[extras.LoginAttempt])
		if response.Content.WasSuccessful {
			break
		}
	}

	return nil

	//TODO: read from buffer to check if username and password match what is found in the db. If it is found, return True.
	//Otherwise, continue login process

}

func SignUpProcess(sessionState *extras.SessionState) error {
	//TODO: pass in connection, prompt user for user name and password, hash the password, salt it, return a valid connection object.
	//This connection obj will contain user account and the tcp connection to the server, ensuring both client and server knows who it belongs to
	enc := gob.NewEncoder(sessionState.ConnectionWrapper.ConnectionObj)

	fmt.Println("Signup Process")
	fmt.Print("What is your username?  ")
	msgScanner.Scan()
	userName := msgScanner.Text()

	fmt.Print("Write a short description describing yourself:")
	msgScanner.Scan()
	description := msgScanner.Text()
	fmt.Print("Your Password: ")
	msgScanner.Scan()
	password := msgScanner.Text()

	accountCreation := extras.UserAccount{UserName: userName, Password: password, Description: description}

	//send login attempt over network
	err := enc.Encode(extras.Packet[extras.Connection]{Type: "SignUp Attempt", Content: extras.Connection{ConnectionObj: sessionState.ConnectionWrapper.ConnectionObj, Account: accountCreation}})
	if err != nil {
		return err
	}
	return nil
}
