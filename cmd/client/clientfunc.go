package main

import (
	extras "chatserver/structs"
	"fmt"
)

func sendMessage(sessionState *extras.SessionState) *extras.Message {

	sessionState.InputScanner.Scan()
	message := sessionState.InputScanner.Text()
	// 1. Move cursor up 1 line
	fmt.Print("\033[1A")
	// ANSI Escape Sequence to clear the entire current line
	fmt.Print("\033[2K")

	// Carriage return to move cursor to the beginning of the line
	fmt.Print("\r")

	return &extras.Message{Text: message}
}

func CommenceAuthenticationProcess(sessionState *extras.SessionState) error {
	//TODO:: Prompt for login, sign up or guest. Then
	fmt.Print("1. To Log in, press '1' \n 2. To Sign Up, press '2' \n To register as a guest, press any other key ")
	sessionState.InputScanner.Scan()
	input := sessionState.InputScanner.Text()
	if input == "1" {
		return LoginProcess(sessionState)
	}
	if input == "2" {
		return SignUpProcess(sessionState)
	}
	fmt.Println("Logging in as Guest")
	// err := sessionState.Encoder.Encode(extras.LoginAttempt{ConnectionWrapper: extras.Connection{ConnectionObj: sessionState.ConnectionWrapper.ConnectionObj}})
	// if err != nil {
	// 	return err
	// }

	return nil

}

func LoginProcess(sessionState *extras.SessionState) error {
	//TODO: pass in connection, check for user account email and password on server and return a valid Connection Object\

	// Create an encoder and target our buffer

	var response *extras.LoginAttempt
	for !sessionState.AuthenticationProcessDone {

		fmt.Println("Logging in")
		fmt.Print("Your username: ")
		sessionState.InputScanner.Scan()
		userName := sessionState.InputScanner.Text()

		fmt.Print("Your Password: ")
		sessionState.InputScanner.Scan()
		password := sessionState.InputScanner.Text()

		loginAttempt := extras.UserAccount{UserName: userName, Password: password}

		//send login attempt over network
		err := sessionState.Encoder.Encode(&extras.LoginAttempt{Account: loginAttempt})
		if err != nil {
			return err
		}

		serverResponse := <-sessionState.PacketListener
		response = serverResponse.(*extras.LoginAttempt)
		if response.WasSuccessful {
			break
		}
	}
	sessionState.AuthenticationProcessDone = true
	sessionState.ConnectionWrapper.Account = response.Account
	return nil

	//TODO: read from buffer to check if username and password match what is found in the db. If it is found, return True.
	//Otherwise, continue login process

}

func SignUpProcess(sessionState *extras.SessionState) error {
	//TODO: pass in connection, prompt user for user name and password, hash the password, salt it, return a valid connection object.
	//This connection obj will contain user account and the tcp connection to the server, ensuring both client and server knows who it belongs to
	var response *extras.SignUpAttempt
	var accountCreation extras.UserAccount
	for {
		fmt.Println("Signup Process")
		fmt.Print("What is your username? (must be unique)")
		sessionState.InputScanner.Scan()
		userName := sessionState.InputScanner.Text()

		fmt.Print("Write a short description describing yourself:")
		sessionState.InputScanner.Scan()
		description := sessionState.InputScanner.Text()
		fmt.Print("Your Password (must be at least 16 characters and include 1 uppercase, 1 lower case and 1 special character): ")
		sessionState.InputScanner.Scan()
		password := sessionState.InputScanner.Text()

		accountCreation = extras.UserAccount{UserName: userName, Password: password, Description: description}

		//send login attempt over network
		err := sessionState.Encoder.Encode(&extras.SignUpAttempt{Account: accountCreation})
		if err != nil {
			return err
		}
		serverResponse := <-sessionState.PacketListener
		response = serverResponse.(*extras.SignUpAttempt)
		if response.WasSuccessful {
			break
		}
		fmt.Println("Your username and password pair are invalid, please try again")

	}
	sessionState.AuthenticationProcessDone = true
	sessionState.ConnectionWrapper.Account = extras.UserAccount{UserName: accountCreation.UserName, Description: accountCreation.Description}
	return nil
}
