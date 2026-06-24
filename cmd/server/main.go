package main

import (
	"errors"
	"fmt"
	"io"
	"net"
)

func main() {
	ListConnections := ConnectionList{connectionList: map[net.Conn]bool{}}
	msgChannel := make(chan string)
	// newUserChannel := make(chan net.Conn)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}
	fmt.Println("Session started")
	defer listener.Close()
	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Printf("%s", err.Error())
			continue
		}
		ListConnections.addConnection(connection)

		go handleConnections(connection, msgChannel, &ListConnections)
		go SendMessages(msgChannel, &ListConnections)

	}

}

// Receives messages from a msg channel and sends them over.
func SendMessages(msgChannel chan string, ListConnections *ConnectionList) {

	newMsg := <-msgChannel
	for conn := range ListConnections.connectionList {
		fmt.Fprintf(conn, "%s", newMsg)
	}

}

// Listens for new messages coming from a particular client
func handleConnections(sender net.Conn, msgChannel chan string, connectionList *ConnectionList) {

	for {
		buffer := make([]byte, 1024)
		bytesRead, err := sender.Read(buffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("Connection closed gracefully by the remote peer.")
			} else {
				fmt.Printf("Connection closed abruptly or failed with error: %v\n", err)
			}
			connectionList.removeConnection(sender)
			return
		}
		msg := buffer[:bytesRead]
		// for _, conn := range connectionList {
		// 	fmt.Fprintf(conn, "%v\n", string(msg))
		// }
		msgChannel <- string(msg)
		fmt.Printf("%v \n", string(msg))

	}

}
