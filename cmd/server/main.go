package main

import (
	"bytes"
	extras "chatserver/structs"
	"encoding/gob"

	"golang.org/x/crypto/argon2"

	"database/sql"
	"errors"
	"fmt"
	"io"
	"net"
)

func main() {

	listConnections := extras.ConnectionList{Connections: map[net.Conn]bool{}}
	//msgChannel := make(chan string)
	msgChannel := make(chan extras.Packet[any])
	// newUserChannel := make(chan net.Conn)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}
	dbConnection, err := sql.Open("sqlite", "ChatServerDB")
	if err != nil {
		fmt.Printf("Database Error: %s", err.Error())
		return
	}

	//serverState := extras.Server{Listener: listener, MsgChannel: msgChannel, ListConnections: &listConnections, Database: dbConnection}
	serverState := extras.Server{Listener: listener, MsgChannel: msgChannel, ListConnections: &listConnections, Database: dbConnection}

	fmt.Println("Session started")

	defer listener.Close()
	defer dbConnection.Close()
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
		//A new goroutine is started for the specific connection. This connection constantly reads the connection for messages sent
		go handleConnections(connection, &serverState)

	}

}

// Receives messages from a msg channel and sends them over.
func SendMessages(serverState *extras.Server) {
	//loops over values in the channel until the channel is closed
	for newMsg := range serverState.MsgChannel {
		if newMsg.Type == "Login Attempt" {
			var accountDetails extras.Packet[extras.Connection]
			accountDetails = newMsg.Content.(extras.Packet[extras.Connection])
			row := serverState.Database.QueryRow("SELECT * FROM Users WHERE username = ?", accountDetails.Content.Account.UserName)

			results := extras.UserAccount{}

			err := row.Scan(&results.Id, &results.UserName, &results.Password, &results.Description, &results.Password, &results.Salt)
			if err != nil {

			}

			attemptHash := argon2.IDKey([]byte(accountDetails.Content.Account.Password), []byte(results.Salt), extras.Time, extras.Memory, extras.Threads, extras.KeyLength)
			if string(attemptHash) == results.Password {
				//Send success message

			} else {
				//Send Failure message
			}

		}
		if newMsg.Type == "BroadCast" {
			serverState.ListConnections.Key.Lock()

			for conn := range serverState.ListConnections.Connections {
				fmt.Fprintf(conn, "%s", newMsg.Content)
			}

			serverState.ListConnections.Key.Unlock()
		}

	}
}

// Listens for new messages coming from a particular client
func handleConnections(sender net.Conn, serverState *extras.Server) {

	buffer := make([]byte, 4096)
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
		rawData := buffer[:bytesRead]
		readBuffer := bytes.NewBuffer(rawData)
		decoder := gob.NewDecoder(readBuffer)
		var decodedPacket extras.Packet[any]
		err = decoder.Decode(&decodedPacket)
		if err != nil {
			fmt.Printf("Error decoding a packet: %v", err)

		}
		//Havent pushed this yet, trying to make code more modular so that chatserver takes more than just text.
		//This will be good for authentication or other types of packets sent over the chatserver other than direct messages
		serverState.MsgChannel <- decodedPacket

	}

}
