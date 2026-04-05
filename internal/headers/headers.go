package headers

import (
	"fmt"
	"strings"
)

type Headers map[string]string

const CRLF = "\r\n"

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	dataString := string(data)

	dataStringSplit := strings.Split(dataString, CRLF)
	if len(dataStringSplit) == 1 {
		return 0, false, nil
	}
	if dataStringSplit[0] == "" {
		return len(CRLF), true, nil

	}

	header := dataStringSplit[0]
	headerKeyValue := strings.Split(header, ":")
	if strings.Contains(headerKeyValue[0], " ") {
		return 0, false, fmt.Errorf("invalid field-name: %s", headerKeyValue[0])
	}
	headerValue := strings.Join(headerKeyValue[1:], ":")
	headerValue = strings.TrimSpace(headerValue)

	h[headerKeyValue[0]] = headerValue

	return len(dataStringSplit[0]) + len(CRLF), false, nil
}
