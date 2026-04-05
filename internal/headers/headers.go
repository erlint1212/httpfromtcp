package headers

import (
	"fmt"
	"strings"
)

type Headers map[string]string

const CRLF = "\r\n"
const bufferSize = 8

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	data_string := string(data)

	data_string_split := strings.Split(data_string, CRLF)
	if len(data_string_split) == 1 {
		return 0, false, nil
	}
	if data_string_split[0][:len(CRLF)-1] == CRLF {
		return len(data_string_split[0]) + len(CRLF), true, nil

	} 

	header := data_string_split[0]
	header_key_value := strings.Split(header, ":")
	if strings.TrimSpace(header_key_value[0]) != header_key_value[0] {
		return 0, false, fmt.Errorf("Whitespace found in field-name: %s", header_key_value[0])
	}
	header_value := strings.Join(header_key_value[1:], ":")
	header_value = strings.TrimSpace(header_value)

	h[header_key_value[0]] = header_value

	return len(data_string_split[0]) + len(CRLF), false, nil
}
