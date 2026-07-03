package main

import (
	"bytes"
	extras "chatserver/structs"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"net"
)

func receiveMessages(sessionState *extras.SessionState) {
	for {

		buffer := make([]byte, 4096)
		bytesRead, err := sessionState.ConnectionWrapper.ConnectionObj.Read(buffer)
		//if server goes down
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("Connection closed gracefully by the remote peer.")
			} else {
				fmt.Printf("Connection closed abruptly or failed with error: %v\n", err)
			}
			return
		}
		// msg := buffer[:bytesRead]
		// // 1. Move cursor up 1 line
		// //fmt.Println("\033[2A")
		// fmt.Print("\033[1B")
		// fmt.Printf("%v \n", string(msg))
		rawData := buffer[:bytesRead]
		readBuffer := bytes.NewBuffer(rawData)
		decoder := gob.NewDecoder(readBuffer)
		var decodedPacket extras.Packet[any]
		err = decoder.Decode(&decodedPacket)
		if err != nil {
			fmt.Printf("Error decoding a packet: %v", err)

		}
		//Havent pushed this yet, trying to make code more modular so that chatserver takes more than just text.
		//This will be good for authentication or other types of packets sent over the chatserver other than direct messages
		if !sessionState.AuthenticationProcessDone {
			if decodedPacket.Type != "Login Attempt" && decodedPacket.Type != "SignUp Attempt" {
				continue

			}
		}
		sessionState.PacketListener <- decodedPacket

	}
}

func main() {
	//send over the tcp protocol

	conn, err := net.Dial("tcp", "localhost:8080")
	packetListener := make(chan extras.Packet[any])
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}
	sessionState := extras.SessionState{ConnectionWrapper: extras.Connection{ConnectionObj: conn}, PacketListener: packetListener}
	go processPackets(&sessionState)
	err = CommenceAuthenticationProcess(&sessionState)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}
	// if connectionObj.Account != (extras.UserAccount{}) {
	// 	//TODO: send login details
	// 	//conn.Write()
	// }

	fmt.Println("\n(Type 'exit()' to close this app)")
	fmt.Println("Type in your msg to send stuff to friends!")

	defer conn.Close()
	for {
		go receiveMessages(&sessionState)
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

func processPackets(sessionState *extras.SessionState) {
	for newPacket := range sessionState.PacketListener {
		if newPacket.Type == "Message" && sessionState.AuthenticationProcessDone {
			//print out msg here

		}
	}
}
