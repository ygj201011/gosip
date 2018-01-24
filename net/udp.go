package net

import (
	"fmt"

	"github.com/ghettovoice/gosip/repository"
)

type udpProtocol struct {
	listeners repository.Repository
}

func NewUdpProtocol() {

}

func (udp *udpProtocol) String() string {
	return fmt.Sprintf("%T %p (net: %s, reliable: %s, streamed: %s)", udp, udp, udp.Network(), udp.Reliable(),
		udp.Streamed())
}

func (udp *udpProtocol) Network() string {
	return "udp"
}

func (udp *udpProtocol) Reliable() bool {
	return false
}

func (udp *udpProtocol) Streamed() bool {
	return false
}

func (udp *udpProtocol) Listen(addr string) error {

}

func (udp *udpProtocol) Send(addr string, data []byte) error {

}
