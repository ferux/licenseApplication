package server

import (
	"io"
	"github.com/ferux/dbgMe"
)

//Connection here
type Connection struct {
	srv *Server
	client *ClientInfo
	dbg dbgMe.Debugger
}
//HandleConnection serves the connection with client. It comes already handshaked.
func (c *Connection) HandleConnection() {
	//TODO: Implement
	msg, err := c.client.ReadMessageString()
	if c.handleError(err) {
		return
	}
	if msg != "me-register" {
		c.client.Close()
	}
	_, err = c.client.SendMessageString("true")
	if c.handleError(err) {
		return
	}
	msg, err = c.client.ReadMessageString()
	if c.handleError(err) {
		return
	}
	mclient, err := c.srv.model.SelectLicense(msg)
	if c.handleError(err) {
		return
	}
	c.client.HostID = mclient.HostID
	
	

	panic("not implemented yet")
}

func (c *Connection) handleError(err error) bool {
	if err != nil {
		c.dbg.Ln(err)
		c.client.Close()
		return true
	}
	return false
}