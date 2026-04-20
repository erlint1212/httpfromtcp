package main

import (
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"io"
	"fmt"
)

func handlerSwitch(w io.Writer, req *request.Request) *server.HandlerError {
	fmt.Println("parsed request!", req.RequestLine.RequestTarget)
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		return handleYourProblem(w, req)
	case "/myproblem":
		return handleMyProblem(w, req)
	default:
		return handleDefault(w, req)
	}

}

func handleYourProblem(w io.Writer, req *request.Request) *server.HandlerError {
	handlerError := server.HandlerError{
		StatusCode: response.StatusCodeBadRequest,
		Message: "Your problem is not my problem\n",

	}
	return &handlerError
}

func handleMyProblem(w io.Writer, req *request.Request) *server.HandlerError {
	handlerError := server.HandlerError{
		StatusCode: response.StatusCodeInternalServerError,
		Message: "Woopsie, my bad\n",

	}
	return &handlerError
}

func handleDefault(w io.Writer, req *request.Request) *server.HandlerError {
	w.Write([]byte("All good, frfr\n"))
	return nil
}

