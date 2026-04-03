package request

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	bit_read, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	parsed_read, err := parseRequestLine(bit_read)
	if err != nil {
		return nil, err
	}

	parsed_request := RequestLine{parsed_read[2], parsed_read[1], parsed_read[0]}
	new_request := Request{parsed_request}

	return &new_request, nil
}

func parseRequestLine(bit_read []byte) ([]string, error) {
	read := string(bit_read)
	lines := strings.Split(read, "\r\n")
	request_line := lines[0]
	request_line_split := strings.Split(request_line, " ")
	for _, rune := range(request_line_split[0]) {
		if !unicode.IsUpper(rune) {
			return nil, fmt.Errorf("method is not pure capital alphabetic character: %s", request_line_split[0])
		}
	}
	request_line_split[2] = strings.Split(request_line_split[2], "/")[1]
	if request_line_split[2] != "1.1" {
		return nil, fmt.Errorf("unsuported HTTP version, expected 1.1, got: %s", request_line_split[2])
	}

	return request_line_split, nil
}


