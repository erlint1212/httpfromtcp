package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
)

func handlerSwitch(w *response.Writer, req *request.Request) {
	fmt.Println("parsed request!", req.RequestLine.RequestTarget)
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		handleYourProblem(w, req)
		return
	case "/myproblem":
		handleMyProblem(w, req)
		return
	default:
		handleDefault(w, req)
		return
	}

}

func handleYourProblem(w *response.Writer, req *request.Request) {
	err := w.WriteStatusLine(response.StatusCodeBadRequest)
	if err != nil {
		fmt.Println("[ERROR] failed to write status line: %v", err)
		return
	}

	body := []byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`)

	header := response.GetDefaultHeaders(len(body))
	err = w.WriteHeaders(header)
	_, err = w.WriteBody(body)
}

func handleMyProblem(w *response.Writer, req *request.Request) {
	err := w.WriteStatusLine(response.StatusCodeInternalServerError)
	if err != nil {
		fmt.Println("[ERROR] failed to write status line: %v", err)
		return
	}

	body := []byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`)
	header := response.GetDefaultHeaders(len(body))
	err = w.WriteHeaders(header)
	_, err = w.WriteBody(body)
}

func handleDefault(w *response.Writer, req *request.Request) {
	err := w.WriteStatusLine(response.StatusCodeOK)
	if err != nil {
		fmt.Println("[ERROR] failed to write status line: %v", err)
		return
	}

	body := []byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)
	header := response.GetDefaultHeaders(len(body))
	err = w.WriteHeaders(header)
	_, err = w.WriteBody(body)
}
