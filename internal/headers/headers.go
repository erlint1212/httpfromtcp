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

func IsValidToken(s string) bool {
	if len(s) == 0 {
		return false
	}

	for _, c := range s {
		if isAlphaNum(byte(c)) || isSpecialTchar(byte(c)) {
			continue
		}
		return false
	}

	return true
}

func isAlphaNum(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}

func isSpecialTchar(c byte) bool {
	const specials = "!#$%&'*+-.^_`|~"
	return strings.IndexByte(specials, c) >= 0
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
	headerKeyValue := strings.SplitN(header, ":", 2)
	if len(headerKeyValue) != 2 {
		return 0, false, fmt.Errorf("invalid header line: %s", header)
	}
	headerKey := headerKeyValue[0]
	if strings.Contains(headerKey, " ") {
		return 0, false, fmt.Errorf("invalid field-name: %s", headerKey)
	}
	if !IsValidToken(headerKey) {
		return 0, false, fmt.Errorf("invalid char in field-name: %s", headerKey)
	}
	headerKey = strings.ToLower(headerKey)
	
	headerValue := headerKeyValue[1]
	headerValue = strings.TrimSpace(headerValue)

	h[headerKey] = headerValue

	return len(dataStringSplit[0]) + len(CRLF), false, nil
}
