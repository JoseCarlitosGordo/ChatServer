package extras

import (
	"database/sql"
	"net"
	"sync"
)

// Configuration Parameters (RFC 9106 Recommendation)
var Time uint32 = 1           // Iterations over memory
var Memory uint32 = 64 * 1024 // 64 MB of RAM
var Threads uint8 = 4         // Parallelism (CPU threads)
var KeyLength uint32 = 32     // Desired output key size

type UserAccount struct {
	Id          int
	UserName    string
	Description string
	Password    string
	Salt        string
}
type ConnectionList struct {
	Key         sync.RWMutex
	Connections map[net.Conn]bool
}

type Packet[T any] struct {
	Type    string
	Content T
}
type Connection struct {
	ConnectionObj net.Conn
	Account       UserAccount
}
type Server struct {
	Listener net.Listener
	//MsgChannel      chan string
	MsgChannel      chan Packet[any]
	ListConnections *ConnectionList
	Database        *sql.DB
}

func (c *ConnectionList) AddConnection(connectionToAdd net.Conn) {
	c.Key.Lock()
	defer c.Key.Unlock()
	c.Connections[connectionToAdd] = true
}

func (c *ConnectionList) RemoveConnection(connectionToRemove net.Conn) {
	c.Key.Lock()
	defer c.Key.Unlock()
	delete(c.Connections, connectionToRemove)
}

func (c *Server) addUser(account UserAccount) {

}
