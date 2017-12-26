// core package defines common library types and methods.
package core

import (
	"strings"

	"github.com/ghettovoice/gosip/util"
)

const (
	DefaultHost     = "127.0.0.1"
	DefaultProtocol = "TCP"

	DefaultUdpPort Port = 5060
	DefaultTcpPort Port = 5060
	DefaultTlsPort Port = 5061
)

// Port number
type Port uint16

func (port *Port) Clone() *Port {
	if port == nil {
		return nil
	}
	newPort := *port
	return &newPort
}

// String wrapper
type MaybeString interface {
	String() string
}

type String struct {
	Str string
}

func (str String) String() string {
	return str.Str
}

// Cancellable can be canceled through cancel method
type Cancellable interface {
	Cancel()
}

type Deferred interface {
	Done() <-chan struct{}
}

const RFC3261BranchMagicCookie = "z9hG4bK"

// GenerateBranch returns random unique branch ID.
func GenerateBranch() string {
	// TODO: use UUID
	return strings.Join([]string{
		RFC3261BranchMagicCookie,
		util.RandStr(16),
	}, "")
}

// DefaultPort returns protocol default port by network.
func DefaultPort(protocol string) Port {
	switch strings.ToLower(protocol) {
	case "tls":
		return DefaultTlsPort
	case "tcp":
		return DefaultTcpPort
	case "udp":
		return DefaultUdpPort
	default:
		return DefaultTcpPort
	}
}
