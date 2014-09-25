package stpclient

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

type stpClient struct {
	// shared
	mu   sync.Mutex
	conn net.Conn
	err  error

	address string

	// Read
	readTimeout time.Duration
	br          *bufio.Reader

	//Write
	writeTimeout time.Duration
	bw           *bufio.Writer
}

func Dial(network, address string) (STPClient, error) {
	c, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return NewConn(c, 0, 0), nil
}

// DialTimeout acts like Dial but takes timeouts for establishing the
// connection to the server, writing a command and reading a reply.
func DialTimeout(network, address string, connectTimeout, readTimeout, writeTimeout time.Duration) (STPClient, error) {
	var c net.Conn
	var err error
	if connectTimeout > 0 {
		c, err = net.DialTimeout(network, address, connectTimeout)
	} else {
		c, err = net.Dial(network, address)
	}
	if err != nil {
		return nil, err
	}
	return NewConn(c, readTimeout, writeTimeout), nil
}

func NewConn(netConn net.Conn, readTimeout, writeTimeout time.Duration) STPClient {
	return &stpClient{
		conn:         netConn,
		bw:           bufio.NewWriter(netConn),
		br:           bufio.NewReader(netConn),
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
	}
}

func (c *stpClient) Close() error {
	c.mu.Lock()
	err := c.err
	if c.err == nil {
		c.err = errors.New("stpClient: closed")
		err = c.conn.Close()
	}
	c.mu.Unlock()
	return err
}

func (c *stpClient) fatal(err error) error {
	c.mu.Lock()
	if c.err == nil {
		c.err = err
		// Close connection to force errors on subsequent calls and to unblock
		// other reader or writer.
		c.conn.Close()
	}
	c.mu.Unlock()
	return err
}

func (c *stpClient) Err() error {
	c.mu.Lock()
	err := c.err
	c.mu.Unlock()
	return err
}

func (c *stpClient) String() string {
	return fmt.Sprintf("stpClient<address=%s>", c.address)
}

func (c *stpClient) writeString(s string) error {
	_, err := c.bw.WriteString(s)
	return err
}

func (c *stpClient) writeBytes(p []byte) error {
	_, err := c.bw.Write(p)
	return err
}

func (c *stpClient) Send(p []byte) error {
	if c.writeTimeout != 0 {
		c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	}
	if err := c.writeBytes(p); err != nil {
		return c.fatal(err)
	}
	return nil
}

// Read response from socket
func (c *stpClient) Receive() (resp []string, err error) {
	resp = make([]string, 0, 5)
	for {
		length, err := c.readLine()
		if err != nil {
			return nil, err
		}
		if length == nil {
			break
		}
		data, err := c.readLine()
		if err != nil {
			return nil, err
		}
		s := string(data)
		resp = append(resp, s)
	}
	return resp, err
}

func (c *stpClient) readLine() ([]byte, error) {
	p, err := c.br.ReadSlice('\n')
	if err == bufio.ErrBufferFull {
		return nil, errors.New("stpClient: long response line")
	}
	if err != nil {
		return nil, err
	}
	i := len(p) - 2
	if i < 0 || p[i] != '\r' {
		return nil, errors.New("stpClient: bad response line terminator")
	}
	return p[:i], nil
}

func (c *stpClient) Flush() error {
	if c.writeTimeout != 0 {
		c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	}
	if err := c.bw.Flush(); err != nil {
		return c.fatal(err)
	}
	return nil
}

func (c *stpClient) Request(request *STPRequest) ([]string, error) {
	var err error
	// send request
	if request != nil {
		c.Send(request.Serialize())
	}
	// flush write buff
	if err := c.Flush(); err != nil {
		return nil, err
	}
	// read the response
	if c.readTimeout != 0 {
		c.conn.SetReadDeadline(time.Now().Add(c.readTimeout))
	}

	var resp []string
	var e error
	if resp, e = c.Receive(); e != nil {
		return nil, c.fatal(e)
	}
	return resp, err

}
