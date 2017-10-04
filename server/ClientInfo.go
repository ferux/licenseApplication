package server

import (
	"crypto/tls"

	"github.com/ferux/validationService/model"
)

//ClientInfo is a wrapper of model.license. Adds connection reference to work with it.
type ClientInfo struct {
	conn *tls.Conn
	model.License
	closeSubscriber []func(*ClientInfo)
}

//Close the connection to client.
func (c *ClientInfo) Close() {
	for _, f := range c.closeSubscriber {
		f(c)
	}
	c.conn.Close()
}

//Subscribe for close event notification.
func (c *ClientInfo) Subscribe(f func(*ClientInfo)) {
	c.closeSubscriber = append(c.closeSubscriber, f)
}

//SetConnection sets the private connection to specified value.
func (c *ClientInfo) SetConnection(con *tls.Conn) {
	c.conn = con
}

//SendMessage sends byte array to the client.
func (c *ClientInfo) SendMessage(data []byte) (int, error) {
	return c.conn.Write(data)
}

//SendMessageString sends string to the client.
func (c *ClientInfo) SendMessageString(data string) (int, error) {
	return c.SendMessage([]byte(data))
}

//ReadMessage reads message from client.
func (c *ClientInfo) ReadMessage() ([]byte, error) {
	buf := make([]byte, 128)
	n, err := c.conn.Read(buf)
	return buf[:n], err
}

//ReadMessageString reads message from client and converts it to string.
func (c *ClientInfo) ReadMessageString() (string, error) {
	msg, err := c.ReadMessage()
	return string(msg), err
}

//ApplyModel copies all parameters from model to client.
func (c *ClientInfo) ApplyModel(m model.License) {
	c.HostID = m.HostID
	c.Application = m.Application
	c.Expiration = m.Expiration
	c.Status = m.Status
}
