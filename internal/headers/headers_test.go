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
