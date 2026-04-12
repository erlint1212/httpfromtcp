package request

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
	state       RequestState
	Headers     headers.Headers
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
)

func RequestFromReader(reader io.Reader) (*Request, error) {
	const bufferSize = 8
	readToIndex := 0
	newRequest := &Request{}
	newRequest.state = requestStateInitialized
	newRequest.Headers = headers.NewHeaders()

	buffer := make([]byte, bufferSize, bufferSize)

	for newRequest.state != requestStateDone {
		if readToIndex == len(buffer) {
			bufferNew := make([]byte, len(buffer)*2, cap(buffer)*2)
			copy(bufferNew, buffer)
			buffer = bufferNew
		}
		numRead, err := reader.Read(buffer[readToIndex:])
		if err == io.EOF {
			if readToIndex == 0 {
				return nil, fmt.Errorf("unexpected EOF")
			}
		} else if err != nil {
			return nil, err
		}

		readToIndex += numRead

		numParsed, err := newRequest.parse(buffer[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buffer, buffer[numParsed:])

		readToIndex -= numParsed
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
	case requestStateParsingHeaders: {
		bytesConsumed, done, err := r.Headers.UseParse(data)
		if err != nil {
			return 0, err
		}

		if !done {
			return bytesConsumed, nil
		}

		r.state = requestStateDone

		return bytesConsumed, nil
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
