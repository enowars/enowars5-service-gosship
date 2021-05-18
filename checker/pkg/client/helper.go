package client

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/Pallinder/go-randomdata"
)

type User struct {
	Name       string             `json:"name"`
	PrivateKey ed25519.PrivateKey `json:"privateKey"`
}

func GenerateRoomAndPassword() (string, string) {
	room := fmt.Sprintf("%s-%s-%s", randomdata.Adjective(), randomdata.Noun(), randomdata.BoundedDigits(6, 0, 999999))
	pwBuf := make([]byte, 16)
	_, _ = rand.Reader.Read(pwBuf)
	return strings.ToLower(room), hex.EncodeToString(pwBuf)
}

func GenerateNewUser() (*User, error) {
	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%s-%s-%s", randomdata.Adjective(), randomdata.FirstName(randomdata.RandomGender), randomdata.BoundedDigits(6, 0, 999999))
	return &User{
		Name:       strings.ToLower(name),
		PrivateKey: privateKey,
	}, nil
}

func GenerateNoise() string {
	return randomdata.Paragraph()
}

type CloseHandler struct {
	fns []func() error
}

func NewCloseHandler() *CloseHandler {
	return &CloseHandler{
		fns: make([]func() error, 0),
	}
}

func (ch *CloseHandler) Add(fn func() error) {
	ch.fns = append(ch.fns, fn)
}

func (ch *CloseHandler) Execute() {
	for _, f := range ch.fns {
		_ = f()
	}
}

type SessionIO struct {
	Session *ssh.Session
	out     io.Reader
	in      io.WriteCloser
}

func (s *SessionIO) Read(p []byte) (n int, err error) {
	return s.out.Read(p)
}

func (s *SessionIO) Write(p []byte) (n int, err error) {
	return s.in.Write(p)
}

func (s *SessionIO) Close() error {
	return s.in.Close()
}

func CreateSSHSession(ctx context.Context, user, addr string, privateKey ed25519.PrivateKey) (*ssh.Client, *SessionIO, *CloseHandler, error) {
	sshSigner, err := ssh.NewSignerFromSigner(privateKey)
	if err != nil {
		return nil, nil, nil, err
	}

	sshClient, err := GetSSHClient(ctx, user, addr, sshSigner)
	if err != nil {
		return nil, nil, nil, err
	}
	ch := NewCloseHandler()
	ch.Add(sshClient.Close)

	session, err := sshClient.NewSession()
	if err != nil {
		return nil, nil, nil, err
	}
	ch.Add(session.Close)

	err = session.RequestPty("xterm", 40, 80, ssh.TerminalModes{
		ssh.ECHO:  0,
		ssh.IGNCR: 1,
	})
	if err != nil {
		ch.Execute()
		return nil, nil, nil, err
	}

	out, err := session.StdoutPipe()
	if err != nil {
		ch.Execute()
		return nil, nil, nil, err
	}
	in, err := session.StdinPipe()
	if err != nil {
		ch.Execute()
		return nil, nil, nil, err
	}
	sio := &SessionIO{
		Session: session,
		out:     out,
		in:      in,
	}
	err = session.Shell()
	if err != nil {
		ch.Execute()
		return nil, nil, nil, err
	}
	return sshClient, sio, ch, nil
}

func AttachRPCAdminClient(ctx context.Context, client *ssh.Client) (*AdminClient, *CloseHandler, error) {
	rpcChannel, err := OpenRPCChannel(client)
	if err != nil {
		return nil, nil, err
	}
	ch := NewCloseHandler()
	ch.Add(rpcChannel.Close)

	grpcConn, err := CreateNewGRPCClient(ctx, rpcChannel)
	if err != nil {
		ch.Execute()
		return nil, nil, err
	}
	ch.Add(grpcConn.Close)

	adminClient := NewAdminClient(grpcConn)

	_, err = adminClient.Auth()
	if err != nil {
		ch.Execute()
		return nil, nil, err
	}
	return adminClient, ch, nil
}
