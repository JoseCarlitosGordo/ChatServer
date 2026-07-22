package main

import (
	"bufio"
	extras "chatserver/structs"
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

func receiveMessages(sessionState *extras.SessionState) {
	for {

		// msg := buffer[:bytesRead]
		// // 1. Move cursor up 1 line
		// //fmt.Println("\033[2A")
		// fmt.Print("\033[1B")
		// fmt.Printf("%v \n", string(msg))

		var decodedPacket extras.Packet
		err := sessionState.Decoder.Decode(&decodedPacket)
		if err != nil {
			fmt.Printf("Error decoding a packet: %v", err)

		}
		//checks if authentication process isnt done and whether the decoded packet is a message
		if !sessionState.AuthenticationProcessDone {
			_, ok := decodedPacket.(*extras.Message)
			if ok {
				continue
			}
		}
		//fmt.Printf("Message Received: %+v \n", decodedPacket)
		sessionState.PacketListener <- decodedPacket

	}
}

func main() {
	//send over the tcp protocol

	extras.RegisterEncodingDecodingTypes()
	conn, err := net.Dial("tcp", "localhost:8080")
	packetListener := make(chan extras.Packet)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}

	sessionState := extras.SessionState{ConnectionWrapper: extras.Connection{ConnectionObj: conn}, PacketListener: packetListener, Decoder: gob.NewDecoder(conn), Encoder: gob.NewEncoder(conn), InputScanner: bufio.NewScanner(os.Stdin)}
	go receiveMessages(&sessionState)
	err = CommenceAuthenticationProcess(&sessionState)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}
	go processPackets(&sessionState)

	fmt.Println("\n(Type 'exit()' to close this app)")
	fmt.Println("Type in your msg to send stuff to friends!")

	defer conn.Close()
	for {
		msgToSend := sendMessage(&sessionState)
		if msgToSend.Text == "exit()" {
			fmt.Println("Exiting the server....")
			return
		}
		msgToSend.Account = sessionState.ConnectionWrapper.Account
		var packet extras.Packet = msgToSend

		err := sessionState.Encoder.Encode(&packet)
		if err != nil {
			fmt.Printf("The server has stopped working unexpectedly...")
			return
		}

	}
}

func processPackets(sessionState *extras.SessionState) {
	for newPacket := range sessionState.PacketListener {
		newPacket.ProcessClientPacket(sessionState)
	}
}
