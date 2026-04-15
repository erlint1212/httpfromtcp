package server

import (
	"fmt"
	"net"
	"strconv"
	"sync/atomic"
)

type ServerState int

const (
	serverStateShutdown ServerState = iota
	serverStateRunning
)

const network = "tcp"

type Server struct {
	listner net.Listener
	state   ServerState
	isClosed atomic.Bool
}

func newServer() *Server {
	return &Server{}
}

func (s *Server) GetAddr() net.Addr {
    if s.listner == nil {
        return nil
    }
    return s.listner.Addr()
}

func Serve(port int) (srv *Server, err error) {
	srvListner, err := net.Listen(network, ":"+strconv.Itoa(port))
	if err != nil {
		return nil, fmt.Errorf("failed to start listen to connection: %v", err)
	}
	srv = newServer()
	srv.listner = srvListner
	srv.state = serverStateRunning

	go srv.listen()

	return srv, nil
}

func (s *Server) Close() (err error) {
	s.isClosed.Store(true)
	err = s.listner.Close()
	if err != nil {
		s.isClosed.Store(false)
		return fmt.Errorf("failed to shutdnow: %v", err)
	}
	s.state = serverStateShutdown

	return nil
}

func (s *Server) listen() {
	for s.state != serverStateShutdown {
		conn, err := s.listner.Accept()
		if err != nil {
			if s.isClosed.Load() {
				fmt.Printf("Server closed, stopping listner loop and ignoring error: %v\n", err)
				return
			}
			fmt.Printf("[ERROR] error while listening: %v\n", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	const response = "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"\r\n" +
		"Hello World!\r\n"
	
	_, err := conn.Write([]byte(response))
	if err != nil {
		fmt.Printf("[ERROR] failed to write to conn: %v\n", err)
	}

	fmt.Println("Response sent to:", conn.RemoteAddr())
}
