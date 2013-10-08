/**
	wire logging utils
 */
package gorets

import (
	"io"
	"net"
	"time"
)

type WireLogConn struct {
	log io.WriteCloser
	conn net.Conn
	// TODO make use of channels to make this perform well
}

func (w *WireLogConn) Read(b []byte) (n int, err error) {
	read, err := w.conn.Read(b)
	w.log.Write(b[0:read])
	return read, err
}
func (w *WireLogConn) Write(b []byte) (n int, err error) {
	write, err := w.conn.Write(b)
	w.log.Write(b[0:write])
	return write, err
}
func (w *WireLogConn) Close() error {
	return w.conn.Close()
}
func (w *WireLogConn) LocalAddr() net.Addr {
	return w.conn.LocalAddr()
}
func (w *WireLogConn) RemoteAddr() net.Addr {
	return w.conn.RemoteAddr()
}
func (w *WireLogConn) SetDeadline(t time.Time) error {
	return w.conn.SetDeadline(t)
}
func (w *WireLogConn) SetReadDeadline(t time.Time) error {
	return w.conn.SetReadDeadline(t)
}
func (w *WireLogConn) SetWriteDeadline(t time.Time) error {
	return w.conn.SetWriteDeadline(t)
}

type Dialer func(network, addr string) (net.Conn, error)

func WireLog(log io.WriteCloser) Dialer {
	return func(network, addr string) (net.Conn, error) {
		conn, err := net.Dial(network, addr)
		wire := WireLogConn{
			log: log,
			conn: conn,
		}
		return &wire, err
	}
}

