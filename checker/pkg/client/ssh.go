package client

import (
	"log"

	"golang.org/x/crypto/ssh"
)

func GetSSHClient(signer ssh.Signer) (*ssh.Client, error) {
	sshClient, err := ssh.Dial("tcp", "127.0.0.1:2222", &ssh.ClientConfig{
		Config:          ssh.Config{},
		User:            "client",
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		return nil, err
	}
	return sshClient, nil
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

func GetAdminClient() {

}
