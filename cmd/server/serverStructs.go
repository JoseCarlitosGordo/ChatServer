package main

import (
	"net"
	"sync"
)

type UserAccount struct {
}
type ConnectionList struct {
	key            sync.Mutex
	connectionList map[net.Conn]bool
}

type Connection struct {
	connectionObj net.Conn
	account       UserAccount
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
