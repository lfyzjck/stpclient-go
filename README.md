stpclient for Golang
====================

client lib that can be used to communicate with simpletp server


Usage
-----
```
// Dial to server
func Dial(network, address string) (*STPClient, error)

// Dial to server with timeout
func DialTimeout(network, address string, connectTimeout, readTimeout, writeTimeout time.Duration) (*STPClient, error)
```

```
import "stpclient/stpclent"

conn := Dial("tcp", "127.0.0.1:3333")
defer conn.Close()
// TODO

```

Reference
---------

See also: [stpclient-py](https://github.com/dccmx/stpclient-py)
