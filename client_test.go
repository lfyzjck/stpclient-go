package stpclient

import (
	"fmt"
	"testing"
)

import "github.com/stretchr/testify/assert"

func TestDial(t *testing.T) {
	conn, err := Dial("tcp", "localhost:3399")
	defer conn.Close()
	assert.Equal(t, err, nil)
	assert.Equal(t, conn.address, "localhost:3399")
}

func TestSend(t *testing.T) {
	conn, err := Dial("tcp", "localhost:3399")
	defer conn.Close()
	req := NewSTPRequest([]string{"ping"})
	err = conn.Send(req.Serialize())
	if err != nil {
		fmt.Println("send bytes failed: ", err)
		t.FailNow()
	}
	err = conn.Flush()
	if err != nil {
		fmt.Println("flush buffer failed: ", err)
		t.FailNow()
	}
	resp, err := conn.Receive()
	if err != nil {
		fmt.Println("recevice response failed: ", err)
		t.FailNow()
	}
	fmt.Println(resp)
}
