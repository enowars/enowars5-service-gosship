package sshnet

import (
	"io"
	"net"

	"golang.org/x/crypto/ssh"
)

type Listener struct {
	conChan chan *Conn
}

func NewListener() *Listener {
	return &Listener{conChan: make(chan *Conn)}
}

func (l *Listener) Accept() (net.Conn, error) {
	c, ok := <-l.conChan
	if !ok {
		return nil, io.EOF
	}
	return c, nil
}

func (l *Listener) Close() error {
	close(l.conChan)
	return nil
}

func (l *Listener) Addr() net.Addr {
	return NewLocalAddr()
}

func (l *Listener) PushChannel(c ssh.Channel, session string) {
	l.conChan <- NewConn(c, session)
}
