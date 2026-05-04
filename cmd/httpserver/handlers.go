package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"httpfromtcp/internal/headers"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func handlerSwitch(w *response.Writer, req *request.Request) {
	fmt.Println("parsed request!", req.RequestLine.RequestTarget)

	target := req.RequestLine.RequestTarget
	if strings.HasPrefix(target, "/httpbin/") {
		if strings.HasPrefix(target, "/httpbin/html") {
			handleTrailerHttpBin(w, req)
			return
		}
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
	case "/video":
		handleGetVideo(w, req)
		return
	default:
		handleDefault(w, req)
		return
	}

}

func handleGetVideo(w *response.Writer, req *request.Request) {
	const videoPath = "./assets/vim.mp4"

	videoFile, err := os.ReadFile(videoPath)
	if err != nil {
		fmt.Println("[ERROR] failed to read video file \"", videoPath, "\": ", err)
		return 
	}

	err = w.WriteStatusLine(response.StatusCodeOK)
	if err != nil {
		fmt.Println("[ERROR] failed to write status line: ", err)
		return
	}

	header := headers.NewHeaders()

	header["Content-Type"] = "video/mp4"
	header["Content-Length"] =  strconv.Itoa(len(videoFile))

	err = w.WriteHeaders(header)
	if err != nil {
		fmt.Println("[ERROR] failed to write headers: ", err)
		return
	}

	_, err = w.WriteBody(videoFile)
	if err != nil {
		fmt.Println("[ERROR] failed to write body: ", err)
		return
	}

}

func handleTrailerHttpBin(w *response.Writer, req *request.Request) {
	resp, err := http.Get("https://httpbin.org/html")
	if err != nil {
		fmt.Printf("[ERROR] failed to GET response from \"https://httpbin.org/html\": %v\n", err)
		return
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Println("[ERROR] failed to close response body: ", err)
		}
	}()

	err = w.WriteStatusLine(response.StatusCode(resp.StatusCode))
	if err != nil {
		fmt.Println("[ERROR] failed to write status line: ", err)
		return
	}

	header := headers.NewHeaders()

	header["Transfer-Encoding"] = "chunked"
	header["Trailer"] = "X-Content-Sha256, X-Content-Length"

	err = w.WriteHeaders(header)
	if err != nil {
		fmt.Println("[ERROR] failed to write headers: ", err)
		return
	}

	buffer := make([]byte, 1024)
	var completeBody []byte

	for {
		bytesRead, err := resp.Body.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("[ERROR] failed to READ response: ", err)
			return
		}

		completeBody = append(completeBody, buffer[:bytesRead]...)

		_, err = w.WriteChunkedBody(buffer[:bytesRead])
		if err != nil {
			fmt.Println("[ERROR] failed to write buffer ", err)
			return
		}

	}

	_, err = w.WriteChunkedBodyDone()
	if err != nil {
		fmt.Println("[ERROR] failed to write body ", err)
		return
	}

	response_hash := sha256.Sum256(completeBody)

	trailer := headers.NewHeaders()

	trailer["X-Content-SHA256"] = hex.EncodeToString(response_hash[:])
	trailer["X-Content-Length"] = fmt.Sprintf("%d", len(completeBody))

	err = w.WriteTrailers(trailer)
	if err != nil {
		fmt.Println("[ERROR] failed to write trailer: ", err)
		return
	}

}

func handleHttpBin(w *response.Writer, req *request.Request) {
	rest := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")

	resp, err := http.Get(fmt.Sprintf("https://httpbin.org/%s", rest))
	if err != nil {
		fmt.Printf("[ERROR] failed to GET response from \"https://httpbin.org/%s\": %v\n", rest, err)
		return
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Println("[ERROR] failed to close response body: ", err)
		}
	}()

	err = w.WriteStatusLine(response.StatusCode(resp.StatusCode))
	if err != nil {
		fmt.Println("[ERROR] failed to write status line: ", err)
		return
	}

	header := headers.NewHeaders()

	header["Transfer-Encoding"] = "chunked"

	err = w.WriteHeaders(header)
	if err != nil {
		fmt.Println("[ERROR] failed to write headers: ", err)
		return
	}

	body := make([]byte, 1024)
	totalWritten := 0

	for {
		bytesRead, err := resp.Body.Read(body)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("[ERROR] failed to READ response: ", err)
			return
		}

		bytesWritten, err := w.WriteChunkedBody(body[:bytesRead])
		if err != nil {
			fmt.Println("[ERROR] failed to write body ", err)
			return
		}
		totalWritten += bytesWritten

	}

	n, err := w.WriteChunkedBodyDone()
	if err != nil {
		fmt.Println("[ERROR] failed to write body ", err)
		return
	}
	totalWritten += n

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
		fmt.Println("[ERROR] failed to write status line: ", err)
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
		fmt.Println("[ERROR] failed to write status line: ", err)
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
