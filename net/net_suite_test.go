package net_test

import (
	"net"
	"strings"
	"testing"

	"github.com/ghettovoice/gosip/test"
	. "github.com/onsi/ginkgo"
)

var (
	localAddr1 = "127.0.0.1:5060"
)

func TestNet(t *testing.T) {
	test.InitTestSuite(t, "Net Suite")
}

func createStreamClientServer(network string, addr string) (net.Conn, net.Conn) {
	network = strings.ToLower(network)
	ln, err := net.Listen(network, addr)
	if err != nil {
		Fail(err.Error())
	}

	ch := make(chan net.Conn)
	go func() {
		defer ln.Close()
		if server, err := ln.Accept(); err == nil {
			ch <- server
		} else {
			Fail(err.Error())
		}
	}()

	client, err := net.Dial(network, ln.Addr().String())
	if err != nil {
		Fail(err.Error())
	}

	return client, <-ch
}

func createPacketClientServer(network string, addr string) (net.Conn, net.Conn) {
	network = strings.ToLower(network)
	server, err := net.ListenPacket(network, addr)
	if err != nil {
		Fail(err.Error())
	}

	client, err := net.Dial(network, server.LocalAddr().String())
	if err != nil {
		Fail(err.Error())
	}

	return client, server.(net.Conn)
}
