package main

import (
	"net"
)

func main() {
	go msgSender()
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		return
	}

	defer listener.Close()
}
