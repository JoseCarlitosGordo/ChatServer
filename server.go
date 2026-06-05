package main

import (
	"fmt"
	"net"
)

func main() {
	go msgSender()
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		return
	}

	defer listener.Close()
	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Printf("%s", err.Error())
		}
		go handleConnections(connection)
	}

}

func handleConnections(connection net.Conn) {

}
