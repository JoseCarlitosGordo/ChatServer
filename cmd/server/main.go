package main

import (
	extras "chatserver/structs"
	"encoding/gob"

	"database/sql"
	"fmt"
	"net"
)

func main() {
	//TODO: Remap connectionList so that each net.conn matches an extras.Connection Obj
	//This is good for ensuring that encoders can be accessed easily in handleConnections
	listConnections := extras.ConnectionList{Connections: map[net.Conn]*extras.Connection{}}

	msgChannel := make(chan extras.Packet)

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
		listConnections.AddConnection(&extras.Connection{ConnectionObj: connection, Encoder: gob.NewEncoder(connection), Decoder: gob.NewDecoder(connection)})
		//A new goroutine is started for the specific connection. This connection constantly reads the connection for messages sent
		go handleConnections(connection, &serverState)

	}

}

// Receives messages from a msg channel and sends them over.
func SendMessages(serverState *extras.Server) {
	//loops over values in the channel until the channel is closed
	for newMsg := range serverState.MsgChannel {
		newMsg.ProcessServerPacket(serverState)
	}
}

// Listens for new messages coming from a particular client
func handleConnections(sender net.Conn, serverState *extras.Server) {
	currentConnection := serverState.ListConnections.Connections[sender]
	for {

		var decodedPacket extras.Packet
		err := currentConnection.Decoder.Decode(&decodedPacket)
		if err != nil {
			fmt.Printf("Error decoding a packet: %v", err)
		}
		//messages that are decoded are sent to a channel where the contents are processed
		serverState.MsgChannel <- decodedPacket

	}

}
