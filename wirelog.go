/**
wire logging utils
*/
package gorets_client

import (
	"io"
	"net"
)

/** this just makes the return type for the Dialer function reasonable */
type Dialer func(network, addr string) (net.Conn, error)

/** create a net.Dial function based on this log */
func WireLog(log io.WriteCloser) Dialer {
	return func(network, addr string) (net.Conn, error) {
		conn, err := net.Dial(network, addr)
		wire := WireLogConn{
			log:  log,
			Conn: conn,
		}
		return &wire, err
	}
}

// channels might make this perform better, though we'ld have to copy the []byte to do that
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
