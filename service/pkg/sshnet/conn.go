package sshnet

import (
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

type Conn struct {
	ssh.Channel
	local, remote net.Addr
}

func NewConn(channel ssh.Channel, session string) *Conn {
	return &Conn{
		Channel: channel,
		local:   NewLocalAddr(),
		remote:  NewRemoteAddr(session),
	}
}

func (c *Conn) LocalAddr() net.Addr {
	return c.local
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.remote
}

func (c *Conn) SetDeadline(t time.Time) error {
	return nil
}

func (c *Conn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *Conn) SetWriteDeadline(t time.Time) error {
	return nil
}
