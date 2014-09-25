package stpclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSTPRequestSerialize(t *testing.T) {
	input := []string{"toupper", "abcd"}
	output := "7\r\ntoupper\r\n4\r\nabcd\r\n\r\n"
	req := &STPRequest{args: input}
	assert.Equal(t, string(req.Serialize()), output, "should be equal")
}
