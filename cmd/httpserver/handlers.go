package main

import (
	"bytes"
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"io"
	"net/http"
	"strings"
)

func handlerSwitch(w *response.Writer, req *request.Request) {
	fmt.Println("parsed request!", req.RequestLine.RequestTarget)

	target := req.RequestLine.RequestTarget
	if strings.HasPrefix(target, "/httpbin/") {
		handleHttpBin(w, req)
		return
	}

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		handleYourProblem(w, req)
		return
	case "/myproblem":
		handleMyProblem(w, req)
		return
	case "/httpbin":
		handleHttpBin(w, req)
		return
	default:
		handleDefault(w, req)
		return
	}

}

func handleHttpBin(w *response.Writer, req *request.Request) {
	err := w.WriteStatusLine(response.StatusCodeOK)
	if err != nil {
		fmt.Println("[ERROR] failed to write status line: %v", err)
		return
	}

	rest := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")

	resp, err := http.Get(fmt.Sprintf("https://httpbin.org/%s", rest))
	if err != nil {
		fmt.Println("[ERROR] failed to GET response from \"https://httpbin.org/x\": ", err)
		return
	}

	body := make([]byte, 1024)

	// body, err := io.ReadAll(resp.Body)
	_, err = resp.Body.Read(body)
	if err != nil {
		fmt.Println("[ERROR] failed to READ response: ", err)
		return
	}

	header := response.GetDefaultHeaders(len(body))
	err = w.WriteHeaders(header)
	if err != nil {
		fmt.Println("[ERROR] failed to write headers: ", err)
		return
	}
	_, err = w.WriteChunkedBody(body)
	if err != nil {
		fmt.Println("[ERROR] failed to write body ", err)
		return
	}
	_, err = w.WriteChunkedBodyDone()
	if err != nil {
		fmt.Println("[ERROR] failed to write body ", err)
		return
	}


}

func handleYourProblem(w *response.Writer, req *request.Request) {
	err := w.WriteStatusLine(response.StatusCodeBadRequest)
	if err != nil {
		fmt.Println("[ERROR] failed to write status line: ", err)
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
	if err != nil {
		fmt.Println("[ERROR] failed to write headers: ", err)
		return
	}
	_, err = w.WriteBody(body)
	if err != nil {
		fmt.Println("[ERROR] failed to write body: ", err)
		return
	}
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
	if err != nil {
		fmt.Println("[ERROR] failed to write headers: ", err)
		return
	}
	_, err = w.WriteBody(body)
	if err != nil {
		fmt.Println("[ERROR] failed to write body: ", err)
		return
	}
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
	if err != nil {
		fmt.Println("[ERROR] failed to write headers: ", err)
		return
	}
	_, err = w.WriteBody(body)
	if err != nil {
		fmt.Println("[ERROR] failed to write body: ", err)
		return
	}
}
