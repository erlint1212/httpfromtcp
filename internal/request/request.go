package request

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strconv"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
	state       RequestState
	Headers     headers.Headers
	Body        []byte
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type RequestState int

const (
	requestStateInitialized RequestState = iota
	requestStateDone
	requestStateParsingHeaders
	requestStateParsingBody
)

func RequestFromReader(reader io.Reader) (*Request, error) {
	const bufferSize = 8
	readToIndex := 0
	newRequest := &Request{}
	newRequest.state = requestStateInitialized
	newRequest.Headers = headers.NewHeaders()
	newRequest.Body = []byte{}

	buffer := make([]byte, bufferSize, bufferSize)

	for newRequest.state != requestStateDone {
		if readToIndex == len(buffer) {
			bufferNew := make([]byte, len(buffer)*2, cap(buffer)*2)
			copy(bufferNew, buffer)
			buffer = bufferNew
		}
		numRead, err := reader.Read(buffer[readToIndex:])

		readToIndex += numRead

		numParsed, parseErr := newRequest.parse(buffer[:readToIndex])
		if parseErr != nil {
			return nil, parseErr
		}

		copy(buffer, buffer[numParsed:])

		readToIndex -= numParsed

		if err != nil {
			if err == io.EOF {
				if newRequest.state == requestStateDone {
					break
				}
				return nil, fmt.Errorf("incomplete request")
			}
			return nil, err
		}
	}

	return newRequest, nil
}

func parseRequestLine(data []byte) ([]string, int, error) {
	read := string(data)
	lines := strings.Split(read, "\r\n")
	if len(lines) == 1 {
		return nil, 0, nil
	}
	requestLine := lines[0]
	requestLineSplit := strings.Split(requestLine, " ")
	for _, rune := range requestLineSplit[0] {
		if !unicode.IsUpper(rune) {
			return nil, 0, fmt.Errorf("method is not pure capital alphabetic character: %s", requestLineSplit[0])
		}
	}
	requestLineSplit[2] = strings.Split(requestLineSplit[2], "/")[1]
	if requestLineSplit[2] != "1.1" {
		return nil, 0, fmt.Errorf("unsuported HTTP version, expected 1.1, got: %s", requestLineSplit[2])
	}

	return requestLineSplit, len(lines[0]) + 2, nil // +2 for \r\n
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case requestStateInitialized:
		{
			parsedRead, bytesConsumed, err := parseRequestLine(data)
			if err != nil {
				return 0, err
			}
			if bytesConsumed == 0 {
				return 0, nil
			}

			parsedRequest := RequestLine{parsedRead[2], parsedRead[1], parsedRead[0]}

			r.RequestLine = parsedRequest

			r.state = requestStateParsingHeaders

			return bytesConsumed, nil
		}
	case requestStateParsingHeaders:
		{
			bytesConsumed, done, err := r.Headers.UseParse(data)
			if err != nil {
				return 0, err
			}

			if !done {
				return bytesConsumed, nil
			}

			if _, ok := r.Headers.Get("content-length"); !ok {
				r.state = requestStateDone
			} else {
				r.state = requestStateParsingBody
			}

			return bytesConsumed, nil
		}
	case requestStateParsingBody:
		{
			contentLength, ok := r.Headers.Get("content-length")
			if !ok {
				r.state = requestStateDone
				return 0, nil
			}
			contentLengthInt, err := strconv.Atoi(contentLength)
			if err != nil {
				return 0, fmt.Errorf("Could not convert content-length to int: %v", err)
			}
			r.Body = append(r.Body, data...)
			if len(r.Body) > contentLengthInt {
				return 0, fmt.Errorf("content-length is smaller than data length: content-length=%d len(data)=%d", contentLengthInt, len(data))
			} else if len(r.Body) < contentLengthInt {
				return len(data), nil
			}

			r.state = requestStateDone

			return len(data), nil

		}
	case requestStateDone:
		{
			return 0, fmt.Errorf("[ERROR] Trying to read data in  a requestStateDone state")
		}
	default:
		{
			return 0, fmt.Errorf("[ERROR] Unkown state")
		}
	}
}
