package main

import (
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"io"
)

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
		Message: "Your problem is not my problem\n",

	}
	return &handlerError
}

func handleDefault(w io.Writer, req *request.Request) *server.HandlerError {
	handlerError := server.HandlerError{
		StatusCode: response.StatusCodeOK,
		Message: "All good, frfr\n",

	}
	return &handlerError
}

