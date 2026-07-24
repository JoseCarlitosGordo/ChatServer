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
	extras.RegisterEncodingDecodingTypes()
	listConnections := &extras.ConnectionList{Connections: make(map[net.Conn]*extras.Connection)}

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

	serverState := &extras.Server{Listener: listener, MsgChannel: msgChannel, ListConnections: listConnections, Database: dbConnection}
	fmt.Println("Session started")
	defer listener.Close()
	defer dbConnection.Close()

	//Separate goroutine for msg sending. ListConnections is updated automatically.
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
		go processPackets(connection, serverState)
		go handleConnections(newConnection, serverState)

	}

}

// Receives messages from a msg channel and sends them over.
func processPackets(conn net.Conn, serverState *extras.Server) {
	//loops over values in the channel until the channel is closed
	for newMsg := range serverState.MsgChannel {
		newMsg.ProcessServerPacket(conn, serverState)
	}
}

// Listens for new messages coming from a particular client
func handleConnections(sender *extras.Connection, serverState *extras.Server) {
	defer serverState.ListConnections.RemoveConnection(sender)
	for {
		var packaging extras.Package
		var decodedPacket extras.Packet

		err := sender.Decoder.Decode(&packaging)
		// fmt.Printf("Message Received: %+v \n", decodedPacket)
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
		serverState.MsgChannel <- decodedPacket

	}

}
