package main

import (
	extras "chatserver/structs"
	"errors"
	"fmt"
	"io"
	"net"
)

func main() {

	listConnections := extras.ConnectionList{Connections: map[net.Conn]bool{}}
	msgChannel := make(chan string)
	// newUserChannel := make(chan net.Conn)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}
	serverState := extras.Server{Listener: listener, MsgChannel: msgChannel, ListConnections: &listConnections}
	fmt.Println("Session started")
	defer listener.Close()
	//Separate goroutine for msg sending. ListConnections is updated automatically.
	go SendMessages(&serverState)
	for {
		//Listen for a speciic connection. If no other connection comes it just waits forever.
		connection, err := listener.Accept()
		if err != nil {
			fmt.Printf("%s", err.Error())
			continue
		}
		listConnections.AddConnection(connection)
		//A new goroutine is started for the specific connection. This connection constantly
		go handleConnections(connection, &serverState)

	}

}

// Receives messages from a msg channel and sends them over.
func SendMessages(serverState *extras.Server) {
	//loops over values in the channel until the channel is closed
	for newMsg := range serverState.MsgChannel {

		serverState.ListConnections.Key.Lock()

		for conn := range serverState.ListConnections.Connections {
			fmt.Fprintf(conn, "%s", newMsg)
		}
		serverState.ListConnections.Key.Unlock()
	}
}

// Listens for new messages coming from a particular client
func handleConnections(sender net.Conn, serverState *extras.Server) {

	buffer := make([]byte, 1024)
	for {

		bytesRead, err := sender.Read(buffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("Connection closed gracefully by the remote peer.")
			} else {
				fmt.Printf("Connection closed abruptly or failed with error: %v\n", err)
			}
			serverState.ListConnections.RemoveConnection(sender)
			return
		}
		msg := buffer[:bytesRead]
		serverState.MsgChannel <- string(msg)

	}

}
