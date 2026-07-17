package extras

import (
	"bytes"
	"crypto/rand"
	"database/sql"
	"encoding/gob"
	"errors"
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
	ProcessServerPacket(conn net.Conn, serverState *Server)
	ProcessClientPacket(sessionState *SessionState)
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
	WasSuccessful bool
	Account       UserAccount
}

type SignUpAttempt struct {
	WasSuccessful bool
	Account       UserAccount
}
type Message struct {
	Text    string
	Account UserAccount
}
type Connection struct {
	ConnectionObj net.Conn
	Decoder       *gob.Decoder
	Encoder       *gob.Encoder
	Account       UserAccount
}

type ConnectionRemover struct {
	ConnectionToRemove net.Conn
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
func (m *Message) ProcessServerPacket(conn net.Conn, serverState *Server) {
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
func (m *Message) ProcessClientPacket(sessionState *SessionState) {
	if m.Account == (UserAccount{}) {
		fmt.Printf("Guest: %s", m.Text)

	} else {
		fmt.Printf("%s: %s", m.Account.UserName, m.Text)
	}

}
func (l *LoginAttempt) ProcessServerPacket(conn net.Conn, serverState *Server) {
	//guest account edge case

	row := serverState.Database.QueryRow("SELECT * FROM Users WHERE username = ?", l.Account.UserName)

	results := UserAccount{}

	err := row.Scan(&results.Id, &results.UserName, &results.Password, &results.Description, &results.Password, &results.Salt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Printf("Username does not exist")
		} else {
			fmt.Printf("Error processing login attempt: %v", err)
		}
		serverState.ListConnections.Connections[conn].Encoder.Encode(&LoginAttempt{WasSuccessful: false})
		return
	}

	attemptHash := argon2.IDKey([]byte(l.Account.Password), []byte(results.Salt), Time, Memory, Threads, KeyLength)

	if string(attemptHash) == results.Password {
		//Send success message
		serverState.ListConnections.Connections[conn].Account = results
		var attempt Packet = &LoginAttempt{WasSuccessful: true, Account: UserAccount{UserName: results.UserName, Description: results.Description}}
		serverState.ListConnections.Connections[conn].Encoder.Encode(&attempt)

	} else {

		serverState.ListConnections.Connections[conn].Encoder.Encode(&LoginAttempt{WasSuccessful: false})
	}

}

func (su *SignUpAttempt) ProcessServerPacket(conn net.Conn, serverState *Server) {

	row := serverState.Database.QueryRow("SELECT 1 FROM Users WHERE username = ? LIMIT 1", su.Account.UserName)
	var output int
	var attempt Packet

	err := row.Scan(&output)
	if err == nil {
		attempt = &SignUpAttempt{WasSuccessful: false}
		serverState.ListConnections.Connections[conn].Encoder.Encode(&attempt)
		return
	}
	fmt.Println("Don't think it got past this point")
	if !errors.Is(err, sql.ErrNoRows) {
		fmt.Printf("Error found while decoding sign up attempt: %s", err)
		attempt = &SignUpAttempt{WasSuccessful: false}
		serverState.ListConnections.Connections[conn].Encoder.Encode(&attempt)
		return
	}

	if !userNameAndPasswordAreValid(su.Account.Password) {
		fmt.Print("LET ME OUT")
		//send back error msg and return False
		attempt = &SignUpAttempt{WasSuccessful: false}
		serverState.ListConnections.Connections[conn].Encoder.Encode(&attempt)
		return

	}
	fmt.Println("Don't think it got past this point")

	var salt []byte
	rand.Read(salt)
	HashedPassword := argon2.IDKey([]byte(su.Account.Password), salt, Time, Memory, Threads, KeyLength)
	serverState.Database.Exec("INSERT INTO Users(username, description, hashedpassword, salt) Values (?, ?, ?, ?)", su.Account.UserName, su.Account.Description, HashedPassword, salt)
	attempt = &SignUpAttempt{WasSuccessful: true, Account: UserAccount{UserName: su.Account.UserName, Description: su.Account.Description}}
	serverState.ListConnections.Connections[conn].Encoder.Encode(&attempt)

}
func (su *SignUpAttempt) ProcessClientPacket(sessionState *SessionState) {
	if su.WasSuccessful {
		sessionState.AuthenticationProcessDone = true
		sessionState.ConnectionWrapper.Account = UserAccount{UserName: su.Account.UserName, Description: su.Account.Description}

	}

}
func (l *LoginAttempt) ProcessClientPacket(sessionState *SessionState) {
	if l.WasSuccessful {
		sessionState.AuthenticationProcessDone = true
		sessionState.ConnectionWrapper.Account = UserAccount{UserName: l.Account.UserName, Description: l.Account.Description}

	}

}
