package netwrap

import (
	"net"

	"github.com/pkg/errors"
)

func WrapListener(l net.Listener, f func(net.Conn) (net.Conn, error)) net.Listener {
	return &listener{
		Listener: l,
		f:        f,
	}
}

type listener struct {
	net.Listener

	f func(net.Conn) (net.Conn, error)
}

func (l *listener) Accept() (net.Conn, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, errors.Wrap(err, "error accepting underlying connection")
	}

	c, err = l.f(c)
	return c, errors.Wrap(err, "error in listerner wrapping function")
}
