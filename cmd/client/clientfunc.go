package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func sendMessage() string {
	msgScanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("(Type 'exit()' to close this app) Your Message: ")
	msgScanner.Scan()
	message := msgScanner.Text()
	return message
}

func msgSender() {
	//send over the tcp protocol
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	defer conn.Close()
	for {

		msgToSend := sendMessage()
		if msgToSend == "exit()" {
			return
		}
		fmt.Printf("%s", msgToSend)
		fmt.Fprintf(conn, "%s", msgToSend)

	}

}
