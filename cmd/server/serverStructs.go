package main

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
	key            sync.RWMutex
	connectionList map[net.Conn]bool
}

type Connection struct {
	connectionObj net.Conn
	account       UserAccount
}
type Server struct {
	listener        net.Listener
	msgChannel      chan string
	ListConnections ConnectionList
	database        sql.DB
}

func (c *ConnectionList) addConnection(connectionToAdd net.Conn) {
	c.key.Lock()
	defer c.key.Unlock()
	c.connectionList[connectionToAdd] = true
}

func (c *ConnectionList) removeConnection(connectionToRemove net.Conn) {
	c.key.Lock()
	defer c.key.Unlock()
	delete(c.connectionList, connectionToRemove)
}

func (c *Server) addUser(account UserAccount) {

}
