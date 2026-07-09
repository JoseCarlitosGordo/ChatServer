package extras

import (
	"bytes"
	"crypto/rand"
	"database/sql"
	"encoding/gob"
	"fmt"
	"net"
	"sync"

	"golang.org/x/crypto/argon2"
)

// Configuration Parameters (RFC 9106 Recommendation)
var Time uint32 = 1           // Iterations over memory
var Memory uint32 = 64 * 1024 // 64 MB of RAM
var Threads uint8 = 4         // Parallelism (CPU threads)
var KeyLength uint32 = 32     // Desired output key size

type Packet interface {
	//Packet functionality for chat server (broadcasting msgs, login processes, etc.)
	ProcessServerPacket(*Server)
	ProcessClientPacket()
}
type UserAccount struct {
	Id          int
	UserName    string
	Description string
	Password    string
	Salt        string
}
type ConnectionList struct {
	Key         sync.RWMutex
	Connections map[net.Conn]*Connection
}

type LoginAttempt struct {
	WasSuccessful     bool
	ConnectionWrapper Connection
}

type SignUpAttempt struct {
	WasSuccessful     bool
	ConnectionWrapper Connection
}
type Message struct {
	Text              string
	ConnectionWrapper Connection
}
type Connection struct {
	ConnectionObj net.Conn
	Decoder       *gob.Decoder
	Encoder       *gob.Encoder
	Account       UserAccount
}
type ConnectionRemover struct {
	ConnectionToRemove Connection
	WasSuccessful      bool
}
type Server struct {
	Listener net.Listener
	//MsgChannel      chan string
	MsgChannel      chan Packet
	ListConnections *ConnectionList

	Database *sql.DB
	Buffer   *bytes.Buffer
}

func (c *ConnectionList) AddConnection(connectionToAdd *Connection) {
	c.Key.Lock()
	defer c.Key.Unlock()
	c.Connections[connectionToAdd.ConnectionObj] = connectionToAdd
}

func (c *ConnectionList) RemoveConnection(connectionToRemove *Connection) {
	c.Key.Lock()
	defer c.Key.Unlock()
	delete(c.Connections, connectionToRemove.ConnectionObj)
}

func (s *Server) addUser(account UserAccount) {

}

func (c *Connection) ProcessServerPacket() {

}

func (c *Connection) ProcessClientPacket() {

}

// Sends messages to every connected user in the chat server
func (m *Message) ProcessServerPacket(serverState *Server) {
	serverState.ListConnections.Key.Lock()
	for conn := range serverState.ListConnections.Connections {
		// go func(c net.Conn) {
		// 	encoder := gob.NewEncoder(c)
		// 	encoder.Encode(m)
		// }(conn.Encoder)
		serverState.ListConnections.Connections[conn].Encoder.Encode(m)

	}
	serverState.ListConnections.Key.Unlock()

}
func (m *Message) ProcessClientPacket() {
	fmt.Printf("%s: %s", m.ConnectionWrapper.Account.UserName, m.Text)

}
func (l *LoginAttempt) ProcessServerPacket(serverState *Server) {
	//guest account edge case
	if l.ConnectionWrapper.Account == (UserAccount{}) {
		return
	}
	row := serverState.Database.QueryRow("SELECT * FROM Users WHERE username = ?", l.ConnectionWrapper.Account.UserName)

	results := UserAccount{}

	err := row.Scan(&results.Id, &results.UserName, &results.Password, &results.Description, &results.Password, &results.Salt)
	if err != nil {
		fmt.Printf("Error processing login attempt: %v", err)
		return
	}

	attemptHash := argon2.IDKey([]byte(l.ConnectionWrapper.Account.Password), []byte(results.Salt), Time, Memory, Threads, KeyLength)

	if string(attemptHash) == results.Password {
		//Send success message
		serverState.ListConnections.Connections[l.ConnectionWrapper.ConnectionObj].Account = results
		serverState.ListConnections.Connections[l.ConnectionWrapper.ConnectionObj].Encoder.Encode(&LoginAttempt{WasSuccessful: true, ConnectionWrapper: Connection{Account: UserAccount{UserName: results.UserName, Description: results.Description}}})

	} else {

		serverState.ListConnections.Connections[l.ConnectionWrapper.ConnectionObj].Encoder.Encode(&LoginAttempt{WasSuccessful: false})
	}

}

func (su *SignUpAttempt) ProcessServerPacket(serverState *Server) {
	row := serverState.Database.QueryRow("SELECT 1 FROM Users WHERE username = ? LIMIT 1", su.ConnectionWrapper.Account.UserName)
	var output int
	err := row.Scan(&output)
	if err != nil {
		fmt.Printf("Error found while decoding sign up attempt: %s", err)
	}
	if !userNameAndPasswordAreValid(output, su.ConnectionWrapper.Account.Password) {
		//send back error msg and return False
		serverState.ListConnections.Connections[su.ConnectionWrapper.ConnectionObj].Encoder.Encode(&SignUpAttempt{WasSuccessful: false})

	}
	var salt []byte
	rand.Read(salt)
	HashedPassword := argon2.IDKey([]byte(su.ConnectionWrapper.Account.Password), salt, Time, Memory, Threads, KeyLength)
	serverState.Database.Exec("INSERT INTO Users(username, description, hashedpassword, salt) Values (?, ?, ?, ?)", su.ConnectionWrapper.Account.UserName, su.ConnectionWrapper.Account.Description, HashedPassword, salt)
	serverState.ListConnections.Connections[su.ConnectionWrapper.ConnectionObj].Encoder.Encode(&LoginAttempt{WasSuccessful: true})

}
func (su *SignUpAttempt) ProcessClientPacket() {

}
func (l *LoginAttempt) ProcessClientPacket() {

}
