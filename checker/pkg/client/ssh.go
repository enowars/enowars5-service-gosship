package client

import (
	"context"
	"log"
	"net"

	"golang.org/x/crypto/ssh"
)

func GetSSHClient(ctx context.Context, user, addr string, signer ssh.Signer) (*ssh.Client, error) {
	fullAddr := addr + ":2222"
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", fullAddr)
	if err != nil {
		return nil, err
	}

	if dl, ok := ctx.Deadline(); ok {
		if err := conn.SetDeadline(dl); err != nil {
			return nil, err
		}
	}

	c, chans, reqs, err := ssh.NewClientConn(conn, fullAddr, &ssh.ClientConfig{
		Config:          ssh.Config{},
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		return nil, err
	}
	return ssh.NewClient(c, chans, reqs), nil
}

func OpenRPCChannel(sshClient *ssh.Client) (ssh.Channel, error) {
	channel, reqs, err := sshClient.OpenChannel("rpc", nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	go ssh.DiscardRequests(reqs)
	return channel, nil
}
