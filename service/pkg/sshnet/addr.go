package sshnet

import "net"

type Addr struct {
	session string
}

func NewRemoteAddr(session string) net.Addr {
	return &Addr{
		session: session,
	}
}

func NewLocalAddr() net.Addr {
	return NewRemoteAddr("local")
}

func (a *Addr) Network() string {
	return "ssh"
}

func (a *Addr) String() string {
	return a.session
}
