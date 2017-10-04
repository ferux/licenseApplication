package server

import (
	"crypto/tls"
	"fmt"
	"sync"

	"github.com/ferux/dbgMe"
	"github.com/ferux/validationService/model"
)

//Server contains config and active connections.
type Server struct {
	connections map[string]*ClientInfo
	dbg         dbgMe.Debugger
	mu          sync.Mutex
	exitc       chan interface{}
	active      int
	model       *model.Model
}

//Run starts server to listen for incoming connections.
func Run(config *Config, dbg dbgMe.Debugger) {
	var srv Server
	srv.dbg = dbg.New("Server")
	srv.exitc = make(chan interface{})
	srv.model = config.Model
	addrStr := fmt.Sprintf("%s:%s", config.IP, config.Port)
	listener, err := tls.Listen("tcp", addrStr, config.ServerConfig)
	defer listener.Close()
	if err != nil {
		panic("Can't start server. Reason:" + err.Error())
	}
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				dbg.F("Can't open connection, error:", err)
			}
			tlsConn, ok := conn.(*tls.Conn)
			if !ok {
				dbg.F("The connection is not secure. Closing...")
				conn.Close()
			}
			if err := tlsConn.Handshake(); err != nil {
				dbg.F("The connection did not handshaked. Closing...")
				conn.Close()
			}
			client := &ClientInfo{
				conn: tlsConn,
			}
			client.Subscribe(srv.eventClientClose)
			connection := &Connection{
				srv:    &srv,
				client: client,
				dbg:    dbg.New(fmt.Sprintf("clientIp/%s", tlsConn.RemoteAddr().String())),
			}
			go connection.HandleConnection()
		}
	}()
	<-srv.exitc
	dbg.P("Got exit signal. Exiting.")
}

func (s *Server) eventClientClose(c *ClientInfo) {
	s.dbg.F("Client (%v) has been disconnected", c.conn.RemoteAddr().String())
	s.RemoveConnection(c)
}

//Kill the server
func (s *Server) Kill() {
	s.exitc <- nil
}

func (s *Server) incActiveConnection() {
	s.mu.Lock()
	s.active++
	s.mu.Unlock()
}
func (s *Server) decActiveConnection() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.active--
}

//RemoveConnection blah
func (s *Server) RemoveConnection(ci *ClientInfo) {
	s.mu.Lock()
	defer s.decActiveConnection()
	defer s.mu.Unlock()
	path := fmt.Sprintf("%s/%s", ci.HostID, ci.Application)
	delete(s.connections, path)
}

//AppendConnection blah
func (s *Server) AppendConnection(ci *ClientInfo) {
	s.mu.Lock()
	defer s.incActiveConnection()
	defer s.mu.Unlock()
	path := fmt.Sprintf("%s/%s", ci.HostID, ci.Application)
	if _, ok := s.connections[path]; ok {
		return
	}
	s.connections[path] = ci
}
