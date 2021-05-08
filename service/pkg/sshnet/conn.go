package sshnet

import (
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

type Conn struct {
	ssh.Channel
}

func (c *Conn) LocalAddr() net.Addr {
	return generateAddr()
}

func (c *Conn) RemoteAddr() net.Addr {
	return generateAddr()
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
