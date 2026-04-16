package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strconv"
	"strings"
)

type StatusCode int

const (
	StatusCodeOK                  StatusCode = 200
	StatusCodeBadRequest          StatusCode = 400
	StatusCodeInternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	statusLine := ""

	switch statusCode {
	case StatusCodeOK:
		{
			statusLine = "HTTP/1.1 200 OK\r\n"
		}
	case StatusCodeBadRequest:
		{
			statusLine = "HTTP/1.1 400 Bad Request\r\n"

		}
	case StatusCodeInternalServerError:
		{
			statusLine = "HTTP/1.1 500 Internal Server Error\r\n"
		}
	default:
		{
			statusLine = fmt.Sprintf("HTTP/1.1 %d\r\n", statusCode)
		}

	}

	_, err := w.Write([]byte(statusLine))
	if err != nil {
		return err
	}
	return nil

}

func GetDefaultHeaders(contentLen int) headers.Headers {
	header := headers.NewHeaders()
	header["Content-Length"] = strconv.Itoa(contentLen)
	header["Connection"] = "close"
	header["Content-Type"] = "text/plain"

	return header

}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	if len(headers) == 0 {
		w.Write([]byte("\r\n"))
		return nil
	}
	var builder strings.Builder
	for k, v := range headers {
		fmt.Fprintf(&builder, "%s: %s\r\n", k, v)
	}
	fmt.Fprintf(&builder, "\r\n")

	headers_string := builder.String()
	_, err := w.Write([]byte(headers_string))

	return err
}
