//go:build !tiny && !notls

package connection

import (
	"crypto/tls"
	"io"
	"net"
	"time"
)

func ConnectTLS(addr string, insecure bool) (io.ReadWriteCloser, error) {
	d := net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	cfg := &tls.Config{
		InsecureSkipVerify: insecure,
	}

	conn, err := tls.DialWithDialer(&d, "tcp", addr, cfg)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
