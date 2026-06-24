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
	//Separate goroutine for msg sending. ListConnections is updated automatically.
	go SendMessages(msgChannel, &ListConnections)
	for {
		//Listen for a speciic connection. If no other connection comes it just waits forever.
		connection, err := listener.Accept()
		if err != nil {
			fmt.Printf("%s", err.Error())
			continue
		}
		ListConnections.addConnection(connection)
		//A new goroutine is started for the specific connection. This connection constantly
		go handleConnections(connection, msgChannel, &ListConnections)

	}

}

// Receives messages from a msg channel and sends them over.
func SendMessages(msgChannel chan string, ListConnections *ConnectionList) {
	//loops over values in the channel until the channel is closed
	for newMsg := range msgChannel {

		ListConnections.key.Lock()

		for conn := range ListConnections.connectionList {
			fmt.Fprintf(conn, "%s", newMsg)
		}
		ListConnections.key.Unlock()
	}
}

// Listens for new messages coming from a particular client
func handleConnections(sender net.Conn, msgChannel chan string, connectionList *ConnectionList) {

	buffer := make([]byte, 1024)
	for {

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
		msgChannel <- string(msg)

	}

}
