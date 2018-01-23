// Package net provides extended network level types and methods.
package net

import (
	"fmt"
	"io"
	"net"
)

// Generic network packet.
type Packet interface {
	String() string
	Data() []byte
	Length() int
	Net() string
	Source() string
	SetSource(src string)
	Destination() string
	SetDestination(dest string)
}

type packet struct {
	data []byte
	net  string
	src  string
	dest string
}

// NewPacket creates new instance of Packet type.
func NewPacket(net string, data []byte) Packet {
	return &packet{
		net:  net,
		data: data,
	}
}

func (pkt *packet) String() string {
	if pkt == nil {
		return "<nil>"
	}

	return fmt.Sprintf("%T %p (net: %s, len: %d, src: %s, dest: %s)", pkt, pkt, pkt.Net(), pkt.Source(),
		pkt.Destination(), pkt.Length())
}

func (pkt *packet) Data() []byte {
	return pkt.data
}

func (pkt *packet) Length() int {
	return len(pkt.data)
}

func (pkt *packet) Net() string {
	return pkt.net
}

func (pkt *packet) Source() string {
	return pkt.src
}

func (pkt *packet) SetSource(src string) {
	pkt.src = src
}

func (pkt *packet) Destination() string {
	return pkt.dest
}

func (pkt *packet) SetDestination(dest string) {
	pkt.dest = dest
}

// Error is a generic network level error.
type Error interface {
	net.Error
	// Network indicates network level errors
	EOF() bool
}

// IsEOFError checks that error is an EOF error.
func IsEOFError(err error) bool {
	if err, ok := err.(Error); ok {
		return err.EOF()
	}
	return err == io.EOF || err == io.ErrClosedPipe
}

// IsTimeoutError checks that error is an timeout error.
func IsTimeoutError(err error) bool {
	if err, ok := err.(Error); ok {
		return err.Timeout()
	}
	if err, ok := err.(net.Error); ok {
		return err.Timeout()
	}
	return false
}

// IsTemporaryError checks that error is an temporary error.
func IsTemporaryError(err error) bool {
	if err, ok := err.(Error); ok {
		return err.Temporary()
	}
	if err, ok := err.(net.Error); ok {
		return err.Temporary()
	}
	return false
}

// PoolError is an pool operation error.
type PoolError struct {
	Err  error
	Op   string
	Pool string
}

func (err *PoolError) Network() bool   { return false }
func (err *PoolError) Timeout() bool   { return false }
func (err *PoolError) Temporary() bool { return false }
func (err *PoolError) Error() string {
	if err == nil {
		return "<nil>"
	}

	s := fmt.Sprintf("%T %s", err, err.Op)
	if err.Pool != "" {
		s += " [" + err.Pool + "]"
	}
	s += ": " + err.Err.Error()

	return s
}
