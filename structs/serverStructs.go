package extras

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
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

// A Packet is any content that is sent from the Server to the Client or Vice Versa
type Packet interface {
	//Packet functionality for chat server (broadcasting msgs, login processes, etc.)
	ProcessServerPacket(conn net.Conn, serverState *Server)
	ProcessClientPacket(sessionState *SessionState)
}
type UserAccount struct {
	Id          int    `json:"id"`
	UserName    string `json:"username"`
	Description string `json:"description"`
	Password    string `json:"password"`
	Salt        string `json:"salt"`
}
type ConnectionList struct {
	Key         sync.RWMutex
	Connections map[net.Conn]*Connection
}
type Package struct {
	LoginAttempt  *LoginAttempt  `json:"loginAttempt,omitempty"`
	SignupAttempt *SignUpAttempt `json:"signupAttempt,omitempty"`
	Message       *Message       `json:"message,omitempty"`
}

func (p *Package) ReconstructPacket() (Packet, error) {
	if p.LoginAttempt != nil {
		return p.LoginAttempt, nil
	}
	if p.SignupAttempt != nil {
		return p.SignupAttempt, nil
	}
	if p.Message != nil {
		return p.Message, nil
	}
	return nil, fmt.Errorf("Smth went wrong while reconstructing a packet")

}

type LoginAttempt struct {
	WasSuccessful bool        `json:"wasSuccessful"`
	Account       UserAccount `json:"account"`
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
		serverState.ListConnections.Connections[conn].Encoder.Encode(&Package{LoginAttempt: &LoginAttempt{WasSuccessful: false}})
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

func (l *LoginAttempt) ProcessClientPacket(sessionState *SessionState) {
	if l.WasSuccessful {
		sessionState.AuthenticationProcessDone = true
		sessionState.ConnectionWrapper.Account = UserAccount{UserName: l.Account.UserName, Description: l.Account.Description}

	}

}

type SignUpAttempt struct {
	WasSuccessful bool        `json:"wasSuccessful"`
	Account       UserAccount `json:"account"`
}

// Checks username and password against security requirements to check if it passes.
// If it does, hash password with a random salt and store in sqlite db
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

// checks if sign up was successful
func (su *SignUpAttempt) ProcessClientPacket(sessionState *SessionState) {
	if su.WasSuccessful {
		sessionState.AuthenticationProcessDone = true
		sessionState.ConnectionWrapper.Account = UserAccount{UserName: su.Account.UserName, Description: su.Account.Description}

	} else {
		fmt.Println("Your username and password pair are invalid, please try again")

	}

}

// Message struct which accounts for the user that sent it and the contents of the msg
type Message struct {
	Text    string      `json:"text"`
	Account UserAccount `json:"account"`
}

// Sends messages to every connected user in the chat server
func (m *Message) ProcessServerPacket(conn net.Conn, serverState *Server) {
	serverState.ListConnections.Key.Lock()
	for conn := range serverState.ListConnections.Connections {
		// go func(c net.Conn) {
		// 	encoder := gob.NewEncoder(c)
		// 	encoder.Encode(m)
		// }(conn.Encoder)
		var packet Packet = m
		serverState.ListConnections.Connections[conn].Encoder.Encode(&packet)

	}
	serverState.ListConnections.Key.Unlock()

}

// received msg is displayed to the terminal
func (m *Message) ProcessClientPacket(sessionState *SessionState) {
	if m.Account == (UserAccount{}) {
		fmt.Printf("Guest: %s", m.Text)

	} else {
		fmt.Printf("%s: %s", m.Account.UserName, m.Text)
	}

}

type Connection struct {
	ConnectionObj net.Conn
	Decoder       *json.Decoder
	Encoder       *json.Encoder
	Account       UserAccount
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

type Server struct {
	Listener net.Listener
	//MsgChannel      chan string
	MsgChannel      chan Packet
	ListConnections *ConnectionList

	Database *sql.DB
}
