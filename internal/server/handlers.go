package server

import (
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"io"
	"net"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w *response.Writer, req *request.Request)

func handleError(statusCode response.StatusCode, err error, conn io.Writer) {
	w := response.NewWriter(conn)
	w.WriteStatusLine(statusCode)
	body := []byte(fmt.Sprintf(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>%v</p>
  </body>
</html>`, err))
	header := response.GetDefaultHeaders(len(body))
	err = w.WriteHeaders(header)
	_, err = w.WriteBody(body)
}

func (s *Server) handle(conn net.Conn, handler Handler) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		handleError(response.StatusCodeBadRequest, err, conn)
		return
	}

	writer := response.NewWriter(conn)
	handler(writer, req)

	fmt.Println("Response sent to:", conn.RemoteAddr())
}
