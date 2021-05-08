package sshnet

import (
	"io"
	"net"

	"golang.org/x/crypto/ssh"
)

type Listener struct {
	conChan chan ssh.Channel
}

func NewListener() *Listener {
	return &Listener{conChan: make(chan ssh.Channel)}
}

func (l *Listener) Accept() (net.Conn, error) {
	c, ok := <-l.conChan
	if !ok {
		return nil, io.EOF
	}
	return &Conn{c}, nil
}

func (l *Listener) Close() error {
	close(l.conChan)
	return nil
}

func (l *Listener) Addr() net.Addr {
	return generateAddr()
}

func (l *Listener) PushChannel(c ssh.Channel) {
	l.conChan <- c
}
