//go:build !tiny

package connection

import (
	"io"
	"net"
	"time"
)

func Connect(addr string) (io.ReadWriteCloser, error) {
	d := net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	return d.Dial("tcp", addr)
}
