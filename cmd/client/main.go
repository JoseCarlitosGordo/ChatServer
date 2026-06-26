package main

import (
	"errors"
	"fmt"
	"io"
	"net"
)

func receiveMessages(connection net.Conn) {
	for {
		buffer := make([]byte, 1024)
		bytesRead, err := connection.Read(buffer)
		//if server goes down
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("Connection closed gracefully by the remote peer.")
			} else {
				fmt.Printf("Connection closed abruptly or failed with error: %v\n", err)
			}
			return
		}
		msg := buffer[:bytesRead]
		// 1. Move cursor up 1 line
		//fmt.Println("\033[2A")
		fmt.Print("\033[1B")
		fmt.Printf("%v \n", string(msg))

	}
}

func main() {
	//send over the tcp protocol
	conn, err := net.Dial("tcp", "localhost:8080")

	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}
	loggedIn, err := CommenceAuthenticationProcess(conn)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}
	if loggedIn {
		//TODO: send login details
		//conn.Write()
	}

	fmt.Println("\n(Type 'exit()' to close this app)")
	fmt.Print("\n Type in your msg to send stuff to friends! \n")

	defer conn.Close()

	for {

		go receiveMessages(conn)

		msgToSend := sendMessage()
		if msgToSend == "exit()" {
			return
		}

		_, err := fmt.Fprintf(conn, "%s", msgToSend)
		if err != nil {
			fmt.Printf("The server has stopped working unexpectedly...")
			return
		}

	}
}
