package headers

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaderLineParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("HoST: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid single header
	headers = NewHeaders()
	s1 := []byte("Set-Person: lane-loves-go\r\n")
	s2 := []byte("Set-Person: prime-loves-zig\r\n")
	s3 := []byte("Set-Person: tj-loves-ocaml\r\n\r\n")
	data = slices.Concat(s1, s2, s3)
	n, done, err = headers.UseParse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "lane-loves-go, prime-loves-zig, tj-loves-ocaml", headers["set-person"])
	assert.Equal(t, len(data), n)
	assert.True(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Unsupported char in header
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Unsupported char in header
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}

func TestNewHeaders(t *testing.T) {
	h := NewHeaders()
	assert.NotNil(t, h)
	assert.Equal(t, 0, len(h))
}

func TestGet(t *testing.T) {
	h := NewHeaders()
	h["content-type"] = "application/json"

	// Exact lowercase match
	val, ok := h.Get("content-type")
	assert.True(t, ok)
	assert.Equal(t, "application/json", val)

	// Mixed case lookup
	val, ok = h.Get("Content-Type")
	assert.True(t, ok)
	assert.Equal(t, "application/json", val)

	// All uppercase lookup
	val, ok = h.Get("CONTENT-TYPE")
	assert.True(t, ok)
	assert.Equal(t, "application/json", val)

	// Missing key
	val, ok = h.Get("Authorization")
	assert.False(t, ok)
	assert.Equal(t, "", val)
}

func TestIsValidToken(t *testing.T) {
	// Valid alphanumeric tokens
	assert.True(t, IsValidToken("Host"))
	assert.True(t, IsValidToken("Accept"))
	assert.True(t, IsValidToken("x123"))

	// Valid tokens with special tchars
	assert.True(t, IsValidToken("Content-Type"))
	assert.True(t, IsValidToken("X-Custom-Header"))
	assert.True(t, IsValidToken("x!#$%&'*+-.^_`|~"))

	// Empty string
	assert.False(t, IsValidToken(""))

	// Invalid characters
	assert.False(t, IsValidToken("Host Name"))
	assert.False(t, IsValidToken("Host:Name"))
	assert.False(t, IsValidToken("Héader"))
	assert.False(t, IsValidToken("Header\t"))
	assert.False(t, IsValidToken("Host("))
	assert.False(t, IsValidToken("a/b"))
	assert.False(t, IsValidToken("a@b"))
	assert.False(t, IsValidToken("a[b"))
	assert.False(t, IsValidToken("a{b"))
	assert.False(t, IsValidToken("a\"b"))
}

func TestParseIncompleteData(t *testing.T) {
	// No CRLF at all — incomplete, should return 0 and not done
	h := NewHeaders()
	data := []byte("Host: localhost")
	n, done, err := h.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}

func TestParseEmptyLineMeansDone(t *testing.T) {
	// Bare CRLF signals end of headers
	h := NewHeaders()
	data := []byte("\r\n")
	n, done, err := h.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 2, n)
	assert.True(t, done)
}

func TestParseColonInValue(t *testing.T) {
	// Value contains colons (e.g. a URL)
	h := NewHeaders()
	data := []byte("Location: https://example.com:8080/path\r\n\r\n")
	n, done, err := h.Parse(data)
	require.NoError(t, err)
	assert.False(t, done)
	assert.Equal(t, 42, n)
	assert.Equal(t, "https://example.com:8080/path", h["location"])
}

func TestParseValueWhitespaceTrimming(t *testing.T) {
	// Leading and trailing whitespace in value should be trimmed
	h := NewHeaders()
	data := []byte("Accept:   text/html   \r\n\r\n")
	n, done, err := h.Parse(data)
	require.NoError(t, err)
	assert.False(t, done)
	assert.True(t, n > 0)
	assert.Equal(t, "text/html", h["accept"])
}

func TestParseEmptyValue(t *testing.T) {
	// Header with empty value after colon
	h := NewHeaders()
	data := []byte("X-Empty:\r\n\r\n")
	n, done, err := h.Parse(data)
	require.NoError(t, err)
	assert.False(t, done)
	assert.True(t, n > 0)
	assert.Equal(t, "", h["x-empty"])
}

func TestParseValueOnlyWhitespace(t *testing.T) {
	// Value is only spaces — should trim to empty
	h := NewHeaders()
	data := []byte("X-Blank:    \r\n\r\n")
	n, done, err := h.Parse(data)
	require.NoError(t, err)
	assert.False(t, done)
	assert.True(t, n > 0)
	assert.Equal(t, "", h["x-blank"])
}

func TestParseMissingColon(t *testing.T) {
	// No colon at all — invalid header
	h := NewHeaders()
	data := []byte("InvalidHeaderLine\r\n\r\n")
	n, done, err := h.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}

func TestParseCaseNormalization(t *testing.T) {
	// Keys should be lowercased
	h := NewHeaders()
	data := []byte("CONTENT-LENGTH: 42\r\n\r\n")
	n, done, err := h.Parse(data)
	require.NoError(t, err)
	assert.False(t, done)
	assert.True(t, n > 0)

	val, ok := h["content-length"]
	assert.True(t, ok)
	assert.Equal(t, "42", val)

	// Original cased key should not exist
	_, ok = h["CONTENT-LENGTH"]
	assert.False(t, ok)
}

func TestParseDuplicateHeadersMerge(t *testing.T) {
	// Two calls to Parse with the same key should comma-merge
	h := NewHeaders()
	data1 := []byte("Accept: text/html\r\n")
	n, done, err := h.Parse(data1)
	require.NoError(t, err)
	assert.False(t, done)
	assert.True(t, n > 0)

	data2 := []byte("Accept: application/json\r\n\r\n")
	n, done, err = h.Parse(data2)
	require.NoError(t, err)
	assert.False(t, done)
	assert.True(t, n > 0)

	assert.Equal(t, "text/html, application/json", h["accept"])
}

func TestParseSpecialTcharInKey(t *testing.T) {
	// Valid special characters in header name
	h := NewHeaders()
	data := []byte("X-My~Header!: value\r\n\r\n")
	n, done, err := h.Parse(data)
	require.NoError(t, err)
	assert.False(t, done)
	assert.True(t, n > 0)
	assert.Equal(t, "value", h["x-my~header!"])
}

func TestUseParseMultipleDifferentHeaders(t *testing.T) {
	// Multiple distinct headers followed by terminator
	h := NewHeaders()
	data := []byte("Host: example.com\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\n")
	n, done, err := h.UseParse(data)
	require.NoError(t, err)
	assert.True(t, done)
	assert.Equal(t, len(data), n)
	assert.Equal(t, "example.com", h["host"])
	assert.Equal(t, "text/plain", h["content-type"])
	assert.Equal(t, "13", h["content-length"])
}

func TestUseParseOnlyTerminator(t *testing.T) {
	// Just the empty line — no headers, immediate done
	h := NewHeaders()
	data := []byte("\r\n")
	n, done, err := h.UseParse(data)
	require.NoError(t, err)
	assert.True(t, done)
	assert.Equal(t, 2, n)
	assert.Equal(t, 0, len(h))
}

func TestUseParseIncomplete(t *testing.T) {
	// Data without a terminating CRLF pair — not done, no error
	h := NewHeaders()
	data := []byte("Host: example.com\r\n")
	n, done, err := h.UseParse(data)
	require.NoError(t, err)
	assert.False(t, done)
	// Should have consumed the one header line
	assert.Equal(t, len("Host: example.com\r\n"), n)
	assert.Equal(t, "example.com", h["host"])
}

func TestUseParseErrorMidstream(t *testing.T) {
	// Valid header, then an invalid one — should error
	h := NewHeaders()
	s1 := []byte("Host: example.com\r\n")
	s2 := []byte("Bad Header: value\r\n\r\n") // space in key
	data := slices.Concat(s1, s2)
	n, done, err := h.UseParse(data)
	require.Error(t, err)
	assert.False(t, done)
	// Should have consumed the first valid header before erroring
	assert.Equal(t, len(s1), n)
}

func TestUseParseEmptyInput(t *testing.T) {
	// Empty byte slice — incomplete, no error, no progress
	h := NewHeaders()
	data := []byte("")
	n, done, err := h.UseParse(data)
	require.NoError(t, err)
	assert.False(t, done)
	assert.Equal(t, 0, n)
}

func TestParseSpaceBeforeColonInKey(t *testing.T) {
	// Space right before the colon is invalid per HTTP spec
	h := NewHeaders()
	data := []byte("Host : value\r\n\r\n")
	_, _, err := h.Parse(data)
	require.Error(t, err)
}

func TestParseTabInKey(t *testing.T) {
	h := NewHeaders()
	data := []byte("Ho\tst: value\r\n\r\n")
	_, _, err := h.Parse(data)
	require.Error(t, err)
}
