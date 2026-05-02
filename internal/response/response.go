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

type WriterState int

const (
	WriteStateLine WriterState = iota
	WriteStateHeaders
	WriteStateBody
)

type Writer struct {
	w          io.Writer
	StatusLine []byte
	state      WriterState
}

func NewWriter(conn io.Writer) *Writer {
	newWriter := &Writer{
		w:     conn,
		state: WriteStateLine,
	}
	return newWriter
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	length := len(p)
	hexSize := fmt.Sprintf("%x\r\n", length)

	var totalWritten int

	n, err := w.WriteBody([]byte(hexSize))
	totalWritten += n
	if err != nil {
		return totalWritten, err
	}

	n, err = w.WriteBody(p)
	totalWritten += n
	if err != nil {
		return totalWritten, err
	}

	n, err = w.WriteBody([]byte("\r\n"))
	totalWritten += n
	if err != nil {
		return totalWritten, err
	}

	return totalWritten, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	n, err := w.WriteBody([]byte("0\r\n\r\n"))
	if err != nil {
		return n, err
	}

	return n, nil
}

func (w *Writer) htmlification(msg string) ([]byte, error) {
	if w.state != WriteStateBody {
		return []byte{}, fmt.Errorf("expected Writer state to be in %v, got %v", WriteStateBody, w.state)
	}
	html_struct := `
	<html>
	  <head>
		<title>%s</title>
	  </head>
	  <body>
		<h1>%s</h1>
		<p>%s</p>
	  </body>
	</html>`
	statusMessage := strings.ReplaceAll(string(w.StatusLine), "\r\n", "")
	statusMessageList := strings.Split(statusMessage, " ")
	return []byte(fmt.Sprintf(html_struct, statusMessageList[1:], statusMessageList[2:], msg)), nil
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != WriteStateLine {
		return fmt.Errorf("expected Writer state to be in %v, got %v", WriteStateLine, w.state)
	}
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

	w.w.Write([]byte(statusLine))
	return nil

}

func GetDefaultHeaders(contentLen int) headers.Headers {
	header := headers.NewHeaders()
	header["Content-Length"] = strconv.Itoa(contentLen)
	header["Connection"] = "close"
	header["Content-Type"] = "text/html"

	return header

}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != WriteStateLine {
		return fmt.Errorf("expected Writer state to be in %v, got %v", WriteStateLine, w.state)
	}
	w.state = WriteStateHeaders
	if len(headers) == 0 {
		w.w.Write([]byte("\r\n"))
		return nil
	}
	var builder strings.Builder
	for k, v := range headers {
		fmt.Fprintf(&builder, "%s: %s\r\n", k, v)
	}
	fmt.Fprintf(&builder, "\r\n")

	headers_string := builder.String()
	w.w.Write([]byte(headers_string))

	w.state = WriteStateBody

	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != WriteStateBody {
		return 0, fmt.Errorf("expected Writer state to be in %v, got %v", WriteStateHeaders, w.state)
	}
	w.state = WriteStateBody
	w.w.Write(p)
	return len(p), nil
}
