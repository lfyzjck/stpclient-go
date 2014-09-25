package stpclient

type Error string

func (err Error) Error() string { return string(err) }

type STPClient interface {
	// Close closes the connection.
	Close() error

	// Err returns a non-nil value if the connection is broken. The returned
	// value is either the first non-nil value returned from the underlying
	// network connection or a protocol parsing error. Applications should
	// close broken connections.
	Err() error

	// Do sends a command to the server and returns the received reply.
	Request(request *STPRequest) (resp []string, err error)

	// Send writes the command to the client's output buffer.
	Send(q []byte) error

	// Flush flushes the output buffer to the Redis server.
	Flush() error

	// Receive receives a single reply from the Redis server
	Receive() (reply []string, err error)
}
