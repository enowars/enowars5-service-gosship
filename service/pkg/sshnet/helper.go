package sshnet

import "net"

func generateAddr() net.Addr {
	ip := net.ParseIP("127.0.0.1")
	return &net.TCPAddr{
		IP:   ip,
		Port: 0,
	}
}
