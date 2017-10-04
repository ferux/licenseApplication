//Package client built-ins into application for license check.
package client

//TODO: Add retry and timeout.
import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/ferux/dbgMe"
	"github.com/ferux/validationService/cert"
	"github.com/ferux/validationService/model"
	"github.com/shirou/gopsutil/host"
)

//Connection is a common struct
type Connection struct {
	Client       model.License
	Server       url.URL
	IsActive     bool
	ReadChan     chan []byte
	StopRefresh  func()
	certificates []tls.Certificate
	errors       []string
	conn         *tls.Conn
	cancel       chan interface{}
	dbg          dbgMe.Debugger
}

//Type of certificate it will use.
//EmbeddedCert: you specify 'CertString' and 'KeyString'. Also you can specify DecryptionSecret if your info is encoded by aes cipher.
//FileCert: loads cert.pem and key.pem (location reads from CertString and KeyString) from storage
const (
	EmbeddedCert int = iota + 1
	FileCert
)

var tdbg dbgMe.Debugger

//TODO: add appName to connection.clientinfo
//New creates new variable of *client.Connection.
func New(server url.URL, cert cert.HandleCertificate, appName string, debugger dbgMe.Debugger) *Connection {
	var connection Connection
	connection.Server = server
	connection.dbg = debugger.New("client")
	dbg := connection.dbg.New("client/New")
	if cert != nil {
		certs, err := cert.GetCertificate()
		if err != nil {
			dbg.P(err)
		}
		connection.certificates = certs
	}
	connection.ReadChan = make(chan []byte, 1)
	connection.fillClientInfo()
	dbg.P("Connection has been created successfully")
	return &connection
}

//RegisterMe opens connection to server, reads unique identifier of PC and sends this info to verification server. Uses only one time on first start.
func (c *Connection) RegisterMe() error {
	dbg := c.dbg.New("client/RegisterClient")
	config := &tls.Config{
		Certificates:       c.certificates,
		InsecureSkipVerify: true,
	}
	tlsConn, err := tls.Dial("tcp", fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port()), config)
	if err != nil {
		return err
	}

	dbg.P("Successfuly openned a connection to %s", tlsConn.RemoteAddr().String())
	c.conn = tlsConn
	c.IsActive = true
	if err := tlsConn.Handshake(); err != nil {
		return err
	}
	err = c.sendToServerString("me-register")
	if err != nil {
		if err == io.EOF {
			return err
		}
		c.Close()
		return err
	}
	msg, err := c.readFromServerString()
	if err != nil {
		if err == io.EOF {
			return err
		}
		c.Close()
		return err
	}
	if msg != "true" {
		c.Close()
		return fmt.Errorf("Can't register to server")
	}
	err = c.sendToServerString(c.Client.HostID)
	byteMsg, err := c.readFromServer()
	if err != nil {
		if err == io.EOF {
			return err
		}
		c.Close()
		return err
	}
	var curClient model.License
	err = json.Unmarshal(byteMsg, &curClient)
	if err != nil {
		return err
	}
	if curClient.HostID != c.Client.HostID {
		c.Close()
		return fmt.Errorf("HostIDs don't match. %v != %v", curClient.HostID, c.Client.HostID)
	}
	c.Client = curClient
	return nil
}

//Close the connection.
func (c *Connection) Close() error {
	return c.conn.Close()
}

//SetRefreshRate returns the channel which notifies about update. If channel sends nil then update has been successful and you can read client info in your connection.client struct.
func (c *Connection) SetRefreshRate(d time.Duration) <-chan error {
	dbg := c.dbg.New("client/SetRefreshRate")
	exitc := make(chan interface{}, 1)
	c.StopRefresh = func() { exitc <- nil }
	licensec := make(chan error)
	go func() {
		for {
			select {
			case <-time.After(d):
				licensec <- c.RefreshClientInfo()
			case <-exitc:
				close(exitc)
				c.StopRefresh = nil
				dbg.P("Refreshing has been stopped")
				break
			}
		}
	}()
	return licensec
}

//RefreshClientInfo tries to update info about client. If so it applies changes to connection.Client struct. In other case returns error.
func (c *Connection) RefreshClientInfo() error {
	var curClient model.License

	err := c.sendToServerString("me-refresh")
	if err != nil {
		if err == io.EOF {
			return err
		}
		c.Close()
		return err
	}
	byteMsg, err := c.readFromServer()
	if err != nil {
		if err == io.EOF {
			return err
		}
		c.Close()
		return err
	}

	err = json.Unmarshal(byteMsg, &curClient)
	if err != nil {
		return err
	}
	if curClient.HostID != c.Client.HostID {
		c.Close()
		return fmt.Errorf("HostIDs don't match. %v != %v", curClient.HostID, c.Client.HostID)
	}
	c.Client = curClient
	return nil
}

func (c *Connection) fillClientInfo() {
	c.Client.HostID = getHostInfo()
	c.Client.Status = model.StatusPending
	c.Client.Expiration = time.Now().Add(time.Minute * 15)
}

func (c *Connection) sendToServerString(data string) error {
	_, err := c.conn.Write([]byte(data))
	return err
}

func (c *Connection) readFromServerString() (string, error) {
	data, err := c.readFromServer()
	return string(data), err
}

func (c *Connection) readFromServer() ([]byte, error) {
	buf := make([]byte, 128)
	n, err := c.conn.Read(buf)
	return buf[:n], err
}

func handleErrors(err error, dbg dbgMe.Debugger) {
	if err != nil {
		dbg.P(err)
	}
}

func getHostInfo() string {
	h, _ := host.Info()
	return h.HostID
}
