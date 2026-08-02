package main

import (
	extras "chatserver/structs"
	"database/sql"
	"encoding/json"
	"fmt"
	"net"

	_ "modernc.org/sqlite"
)

func main() {
	//TODO: Remap connectionList so that each net.conn matches an extras.Connection Obj
	//This is good for ensuring that encoders can be accessed easily in handleConnections

	listConnections := &extras.ConnectionList{Connections: make(map[net.Conn]*extras.Connection)}

	// msgChannel := make(chan extras.Packet, 100)
	msgChannel := make(chan extras.InboundMessage)
	listener, err := net.Listen("tcp", "[202:960c:a87:2609:9ec0:29f8:e21d:690a]:8080")
	if err != nil {
		for {
			fmt.Printf("%s", err.Error())
		}
		return
	}
	dbConnection, err := sql.Open("sqlite", "ChatServerDB")
	if err != nil {
		fmt.Printf("Database Error: %s", err.Error())
		return
	}

	serverState := &extras.Server{Listener: listener, MsgChannel: msgChannel, ListConnections: listConnections, Database: dbConnection}
	fmt.Println("Session started")
	defer listener.Close()
	defer dbConnection.Close()

	//Separate goroutine for msg sending. ListConnections is updated automatically.
	go processPackets(serverState)
	for {
		//Listen for a speciic connection. If no other connection comes it just waits forever.
		connection, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error Accepting Connection: %v", err.Error())
			continue
		}
		newConnection := &extras.Connection{ConnectionObj: connection, Encoder: json.NewEncoder(connection), Decoder: json.NewDecoder(connection)}
		serverState.ListConnections.AddConnection(newConnection)
		//A new goroutine is started for the specific connection. This connection constantly reads the connection for messages sent
		go handleConnections(newConnection, serverState)

	}

}

// Receives messages from a msg channel and sends them over.
func processPackets(serverState *extras.Server) {
	//loops over values in the channel until the channel is closed
	for newMsg := range serverState.MsgChannel {

		newMsg.Packet.ProcessServerPacket(&newMsg.Sender, serverState)

	}
}

// Listens for new messages coming from a particular client
func handleConnections(sender *extras.Connection, serverState *extras.Server) {
	defer serverState.ListConnections.RemoveConnection(sender)
	for {
		var packaging extras.Package
		var decodedPacket extras.Packet

		err := sender.Decoder.Decode(&packaging)

		if err != nil {
			fmt.Printf("Error decoding a packet: %v", err)
			return
		}
		decodedPacket, err = packaging.ReconstructPacket()
		if err != nil {
			fmt.Printf("%v", err.Error())
			return
		}
		//messages that are decoded are sent to a channel where the contents are processed
		//Keep an eye on this code....
		serverState.MsgChannel <- extras.InboundMessage{Packet: decodedPacket, Sender: *sender}

	}

}
