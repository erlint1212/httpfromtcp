package server

import (
	"fmt"
	"httpfromtcp/internal/request"
	"io"
)

type HandlerError struct {
	StatusCode int
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func WriteHandlerError(w io.Writer, handlerError *HandlerError) error {
	statusLine := fmt.Sprintf("HTTP/1.1 %d %s", handlerError.StatusCode, handlerError.Message)
	_, err := w.Write([]byte(statusLine))
	return err
}
