package main

import (
	"bytes"
	extras "chatserver/structs"
	"encoding/gob"
	"fmt"
	"net"
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
		//Havent pushed this yet, trying to make code more modular so that chatserver takes more than just text.
		//This will be good for authentication or other types of packets sent over the chatserver other than direct messages
		if !sessionState.AuthenticationProcessDone {
			_, ok := decodedPacket.(*extras.Message)
			if ok {
				continue
			}
		}
		sessionState.PacketListener <- decodedPacket

	}
}

func main() {
	//send over the tcp protocol
	var buffer bytes.Buffer
	conn, err := net.Dial("tcp", "localhost:8080")
	packetListener := make(chan extras.Packet)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}
	sessionState := extras.SessionState{ConnectionWrapper: extras.Connection{ConnectionObj: conn}, PacketListener: packetListener, Buffer: &buffer, Decoder: gob.NewDecoder(conn), Encoder: gob.NewEncoder(conn)}
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
		if msgToSend.Text == "exit()" {
			fmt.Println("Exiting the server....")
			return
		}
		msgToSend.ConnectionWrapper = sessionState.ConnectionWrapper

		err := sessionState.Encoder.Encode(msgToSend)
		if err != nil {
			fmt.Printf("The server has stopped working unexpectedly...")
			return
		}

	}
}

func processPackets(sessionState *extras.SessionState) {
	for newPacket := range sessionState.PacketListener {
		// if newPacket.Type == "Message" && sessionState.AuthenticationProcessDone {
		// 	//print out msg here
		// switch t := newPacket.(type) {
		// case *extras.Message:

		// }

		// }
		newPacket.ProcessClientPacket()
	}
}
