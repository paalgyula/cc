//go:build tiny

package connection

import (
	"io"

	net "tinygo.org/x/drivers/net"
)

func Connect(addr string) (io.ReadWriteCloser, error) {
	return net.Dial("tcp", addr)
}
