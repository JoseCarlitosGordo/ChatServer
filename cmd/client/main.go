package main

import (
	"fmt"
	"net"
)

func main() {
	//send over the tcp protocol
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}
	defer conn.Close()
	for {
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
