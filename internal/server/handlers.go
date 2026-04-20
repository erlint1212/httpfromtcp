package server

import (
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"io"
	"fmt"
	"net"
	"bytes"

)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func WriteHandlerError(w io.Writer, hErr *HandlerError) error {
    body := []byte(hErr.Message)
    if err := response.WriteStatusLine(w, hErr.StatusCode); err != nil {
        return err
    }
    headers := response.GetDefaultHeaders(len(body))
    if err := response.WriteHeaders(w, headers); err != nil {
        return err
    }
    _, err := w.Write(body)
    return err
}

func (hErr *HandlerError) Write(conn net.Conn) {
	err := WriteHandlerError(conn, hErr)
	if err != nil {
		fmt.Printf("[ERROR] failed to write error to writer: %v\n", err)
	}
}

func (s *Server) handle(conn net.Conn, handler Handler) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Printf("[ERROR] failed to parse request from connection: %v\n", err)
		handlerError := &HandlerError{
			StatusCode: response.StatusCodeBadRequest,
			Message: err.Error(),
		}
		handlerError.Write(conn)
		return
	}

	var buffer bytes.Buffer

	handlerErr := handler(&buffer, req)
	if handlerErr != nil {
		handlerErr.Write(conn)
		return
	}

	b := buffer.Bytes()
	err = response.WriteStatusLine(conn, response.StatusCodeOK)
	if err != nil {
		fmt.Printf("[ERROR] failed to write status line to conn: %v\n", err)
		return
	}

	statusLine := response.GetDefaultHeaders(len(b))
	err = response.WriteHeaders(conn, statusLine)
	if err != nil {
		fmt.Printf("[ERROR] failed to write headers line to conn: %v\n", err)
		return
	}

	_, err = buffer.WriteTo(conn)
	if err != nil {
		fmt.Printf("[ERROR] failed to write buffer to conn: %v\n", err)
		return
	}

	fmt.Println("Response sent to:", conn.RemoteAddr())
}
