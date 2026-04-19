package server

import (
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"io"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func WriteHandlerError(w io.Writer, handlerError *HandlerError) error {
	statusLine := fmt.Sprintf("HTTP/1.1 %d %s", handlerError.StatusCode, handlerError.Message)
	_, err := w.Write([]byte(statusLine))
	return err
}

func handleYourProblem(w io.Writer, req *request.Request) *HandlerError {
	handlerError := HandlerError{
		StatusCode: response.StatusCodeBadRequest,
		Message: "Your problem is not my problem\n",

	}
	return &handlerError
}

func handleMyProblem(w io.Writer, req *request.Request) *HandlerError {
	handlerError := HandlerError{
		StatusCode: response.StatusCodeInternalServerError,
		Message: "Your problem is not my problem\n",

	}
	return &handlerError
}

func handleDefault(w io.Writer, req *request.Request) *HandlerError {
	handlerError := HandlerError{
		StatusCode: response.StatusCodeOK,
		Message: "All good, frfr\n",

	}
	return &handlerError
}
