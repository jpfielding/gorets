package client

import (
	"crypto/tls"
	"io"
	"net"
)

// Dialer just makes the return type for the Dialer function reasonable
type Dialer func(network, addr string) (net.Conn, error)

// WireLog Dial = gorets.WireLog(file, dial)
func WireLog(log io.WriteCloser, dial Dialer) Dialer {
	return func(network, addr string) (net.Conn, error) {
		conn, err := dial(network, addr)
		wire := WireLogConn{
			log:  log,
			Conn: conn,
		}
		return &wire, err
	}
}

// WireLogTLS Transport.DialTLS = gorets.WireLogTLS(file)
func WireLogTLS(log io.WriteCloser) Dialer {
	return func(network, addr string) (net.Conn, error) {
		config := &tls.Config{InsecureSkipVerify: true}
		c, err := tls.Dial(network, addr, config)
		if err != nil {
			return nil, err
		}
		err = c.Handshake()
		wire := WireLogConn{
			log:  log,
			Conn: c,
		}
		if err != nil {
			return &wire, err
		}
		return &wire, c.Handshake()
	}
}

// WireLogConn ....
type WireLogConn struct {
	// embedded
	net.Conn
	// the destination for the split stream
	log io.WriteCloser
}

func (w *WireLogConn) Read(b []byte) (n int, err error) {
	read, err := w.Conn.Read(b)
	w.log.Write(b[0:read])
	return read, err
}
func (w *WireLogConn) Write(b []byte) (n int, err error) {
	write, err := w.Conn.Write(b)
	w.log.Write(b[0:write])
	return write, err
}
