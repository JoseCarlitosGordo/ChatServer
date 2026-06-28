package extras

import (
	"database/sql"
	"net"
	"sync"
)

type UserAccount struct {
	UserName    string
	Description string
	Password    string
	salt        string
}
type ConnectionList struct {
	Key         sync.RWMutex
	Connections map[net.Conn]bool
}

type Connection struct {
	ConnectionObj net.Conn
	Account       UserAccount
}
type Server struct {
	Listener        net.Listener
	MsgChannel      chan string
	ListConnections ConnectionList
	Database        sql.DB
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
