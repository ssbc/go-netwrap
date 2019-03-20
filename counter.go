package netwrap

import (
	"net"
	"sync"

	"github.com/pkg/errors"
)

func CountingWrapper(c net.Conn) (net.Conn, error) {
	return &countingConn{Conn: c}, nil
}

type countingConn struct {
	net.Conn

	lock   sync.Mutex
	closed bool

	tx uint64
	rx uint64
}

func (conn *countingConn) Write(b []byte) (int, error) {
	n, err := conn.Conn.Write(b)
	conn.tx += uint64(n)
	return n, err
}

func (conn *countingConn) Read(b []byte) (int, error) {
	n, err := conn.Conn.Read(b)
	conn.rx += uint64(n)
	return n, err
}

func (conn *countingConn) Close() error {
	conn.lock.Lock()
	defer conn.lock.Unlock()
	if conn.closed {
		return errors.Errorf("countingConn: already closed")
	}
	err := conn.Conn.Close()
	conn.closed = true
	return errors.Wrap(err, "countingConn: failed to close")
}
