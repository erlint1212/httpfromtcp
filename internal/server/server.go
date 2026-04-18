package server

import (
	"bytes"
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
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

func Serve(port int, handler Handler) (srv *Server, err error) {
	srvListner, err := net.Listen(network, ":"+strconv.Itoa(port))
	if err != nil {
		return nil, fmt.Errorf("failed to start listen to connection: %v", err)
	}
	srv = newServer()
	srv.listner = srvListner
	srv.state = serverStateRunning

	go srv.listen(handler)

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

func (s *Server) listen(handler Handler) {
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
		go s.handle(conn, handler)
	}
}

func (s *Server) handle(conn net.Conn, handler Handler) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Printf("[ERROR] failed to parse request from connection: %v\n", err)
		return
	}

	var buffer bytes.Buffer
	handlerErr := handler(&buffer, req)
	err = WriteHandlerError(&buffer, handlerErr)
	if err != nil {
		fmt.Printf("[ERROR] failed to write error to writer: %v\n", err)
		return
	}

	err = response.WriteStatusLine(conn, response.StatusCodeOK)
	if err != nil {
		fmt.Printf("[ERROR] failed to write status line to conn: %v\n", err)
		return
	}

	statusLine := response.GetDefaultHeaders(0)
	err = response.WriteHeaders(conn, statusLine)
	if err != nil {
		fmt.Printf("[ERROR] failed to write headers line to conn: %v\n", err)
		return
	}

	fmt.Println("Response sent to:", conn.RemoteAddr())
}
