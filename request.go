package stpclient

import (
	"bytes"
	"fmt"
)

type STPRequest struct {
	args []string
}

// Serialize STPRequest args to a string
func (req STPRequest) Serialize() []byte {
	var buffer bytes.Buffer

	for i := 0; i < len(req.args); i++ {
		buffer.WriteString(fmt.Sprintf("%d\r\n%s\r\n", len(req.args[i]), req.args[i]))
	}
	buffer.WriteString("\r\n")
	return buffer.Bytes()
}
