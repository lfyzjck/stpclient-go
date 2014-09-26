package stpclient

import (
	"bytes"
	"fmt"
)

type STPRequest struct {
	args []string
}

func NewSTPRequest(args []string) *STPRequest {
	return &STPRequest{
		args: args,
	}
}

// Serialize STPRequest args to a string
func (req STPRequest) Serialize() string {
	var buffer bytes.Buffer

	for i := 0; i < len(req.args); i++ {
		line := fmt.Sprintf("%d\r\n%s\r\n", len(req.args[i]), req.args[i])
		buffer.WriteString(line)
	}
	buffer.WriteString("\r\n")
	return buffer.String()
}
