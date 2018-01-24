package net

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/ghettovoice/gosip/log"
)

var (
	bufSize      uint16 = 65535 - 20 - 8 // IPv4 max size - IPv4 Header size - UDP Header size
	readTimeout         = 30 * time.Second
	writeTimeout        = 30 * time.Second
)

// Conn is an extended net.Conn.
type Conn interface {
	net.Conn
	log.LocalLogger
	Network() string
	Streamed() bool
	String() string
}

// Conn implementation.
type conn struct {
	log.LocalLogger
	baseConn net.Conn
	laddr    net.Addr
	raddr    net.Addr
	streamed bool
	mu       *sync.RWMutex
}

// NewConn creates new connection from the base net.Conn connection.
func NewConn(baseConn net.Conn) Conn {
	var stream bool
	switch baseConn.(type) {
	case net.PacketConn:
		stream = false
	default:
		stream = true
	}

	conn := &conn{
		LocalLogger: log.NewSafeLocalLogger(),
		baseConn:    baseConn,
		laddr:       baseConn.LocalAddr(),
		raddr:       baseConn.RemoteAddr(),
		streamed:    stream,
		mu:          new(sync.RWMutex),
	}
	return conn
}

func (conn *conn) String() string {
	if conn == nil {
		return "<nil>"
	}

	return fmt.Sprintf("%T %p (net: %s, laddr: %v, raddr: %v)", conn, conn, conn.Network(), conn.LocalAddr(),
		conn.RemoteAddr())
}

func (conn *conn) Log() log.Logger {
	// remote addr for net.PacketConn resolved in runtime
	return conn.LocalLogger.Log().WithFields(map[string]interface{}{
		"conn":  conn.String(),
		"raddr": fmt.Sprintf("%v", conn.RemoteAddr()),
	})
}

func (conn *conn) SetLog(logger log.Logger) {
	conn.LocalLogger.SetLog(logger.WithFields(map[string]interface{}{
		"laddr": fmt.Sprintf("%v", conn.LocalAddr()),
		"net":   strings.ToUpper(conn.LocalAddr().Network()),
	}))
}

func (conn *conn) Streamed() bool {
	return conn.streamed
}

func (conn *conn) Network() string {
	return strings.ToUpper(conn.baseConn.LocalAddr().Network())
}

func (conn *conn) Read(buf []byte) (int, error) {
	var (
		num   int
		err   error
		raddr net.Addr
	)

	if err := conn.baseConn.SetReadDeadline(time.Now().Add(readTimeout)); err != nil {
		conn.Log().Warnf("%s failed to set read deadline: %s", conn, err)
	}

	switch baseConn := conn.baseConn.(type) {
	case net.PacketConn: // UDP & ...
		num, raddr, err = baseConn.ReadFrom(buf)
		conn.mu.Lock()
		conn.raddr = raddr
		conn.mu.Unlock()
	default: // net.Conn - TCP, TLS & ...
		num, err = conn.baseConn.Read(buf)
	}

	if err != nil {
		return num, &ConnError{
			err,
			"read",
			conn.Network(),
			fmt.Sprintf("%v", conn.RemoteAddr()),
			fmt.Sprintf("%v", conn.LocalAddr()),
			conn,
		}
	}

	conn.Log().Debugf(
		"%s received %d bytes",
		conn,
		num,
	)

	return num, err
}

func (conn *conn) Write(buf []byte) (int, error) {
	var (
		num int
		err error
	)

	if err := conn.baseConn.SetWriteDeadline(time.Now().Add(writeTimeout)); err != nil {
		conn.Log().Warnf("%s failed to set write deadline: %s", conn, err)
	}

	num, err = conn.baseConn.Write(buf)
	if err != nil {
		return num, &ConnError{
			err,
			"write",
			conn.Network(),
			fmt.Sprintf("%v", conn.RemoteAddr()),
			fmt.Sprintf("%v", conn.LocalAddr()),
			conn,
		}
	}

	conn.Log().Debugf(
		"%s written %d bytes",
		conn,
		num,
	)

	return num, err
}

func (conn *conn) LocalAddr() net.Addr {
	return conn.laddr
}

func (conn *conn) RemoteAddr() net.Addr {
	// we should protect raddr field with mutex,
	// because there is may be DATA RACE with Read method that usually executes
	// in another goroutine
	conn.mu.RLock()
	defer conn.mu.RUnlock()
	return conn.raddr
}

func (conn *conn) Close() error {
	err := conn.baseConn.Close()
	if err != nil {
		return &ConnError{
			err,
			"close",
			conn.Network(),
			"",
			"",
			conn,
		}
	}

	conn.Log().Debugf(
		"%s closed",
		conn,
	)

	return nil
}

func (conn *conn) SetDeadline(t time.Time) error {
	return conn.baseConn.SetDeadline(t)
}

func (conn *conn) SetReadDeadline(t time.Time) error {
	return conn.baseConn.SetReadDeadline(t)
}

func (conn *conn) SetWriteDeadline(t time.Time) error {
	return conn.baseConn.SetWriteDeadline(t)
}

// Connection level error.
type ConnError struct {
	Err    error
	Op     string
	Net    string
	Source string
	Dest   string
	Conn   Conn
}

func (err *ConnError) EOF() bool       { return IsEOFError(err.Err) }
func (err *ConnError) Timeout() bool   { return IsTimeoutError(err.Err) }
func (err *ConnError) Temporary() bool { return IsTemporaryError(err.Err) }
func (err *ConnError) Error() string {
	if err == nil {
		return "<nil>"
	}

	s := fmt.Sprintf("%T", err)
	if err.Conn != nil {
		s += " [" + err.Conn.String() + "]"
	}
	s += " " + err.Op
	if err.Source != "" {
		s += " " + err.Source
	}
	if err.Dest != "" {
		if err.Source != "" {
			s += "->"
		} else {
			s += " "
		}
		s += err.Dest
	}

	s += ": " + err.Err.Error()

	return s
}
