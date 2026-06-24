package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

func main() {
	var key sync.Mutex
	connectionsList := []net.Conn{}
	msgChannel := make(chan string)
	// newUserChannel := make(chan net.Conn)
	go SendMessages(msgChannel)
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
		connectionsList = append(connectionsList, connection)

		go handleConnections(connection, msgChannel, key)

	}

}

// Receives messages from a msg channel and sends them over.
func SendMessages(msgChannel chan string) {

	for {
		// newMsg := <-msgChannel

	}

}

// Listens for new messages coming from a particular client
func handleConnections(sender net.Conn, msgChannel chan string, key sync.Mutex) {

	for {
		buffer := make([]byte, 1024)
		bytesRead, err := sender.Read(buffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("Connection closed gracefully by the remote peer.")
			} else {
				fmt.Printf("Connection closed abruptly or failed with error: %v\n", err)
			}
			key.Lock()

			key.Unlock()
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

//TODO: Make use of channels
