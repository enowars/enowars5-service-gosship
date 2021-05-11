package main

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"gosship/pkg/rpc/admin"
	"gosship/pkg/sshnet"
	"log"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"google.golang.org/grpc"
)

var privateKey ed25519.PrivateKey

func init() {
	privateKeyRaw, err := hex.DecodeString("215c8787c1b079be149db3da5e297a9b39ff008dee69b2b9f115d51d4547664580de0e58c0842f83cf95f9772a5a13c167dd4c0e3fd02913076d16df828fbbb2")
	if err != nil {
		panic(err)
	}
	privateKey = privateKeyRaw
	log.Println("public key:", hex.EncodeToString(privateKey[32:]))
}

type AdminClient struct {
	svc          admin.AdminServiceClient
	sessionToken string
}

func NewAdminClient(svc admin.AdminServiceClient) *AdminClient {
	return &AdminClient{svc: svc}
}

func (a *AdminClient) Auth() (string, error) {
	authChallenge, err := a.svc.GetAuthChallenge(context.Background(), &admin.GetAuthChallenge_Request{})
	if err != nil {
		return "", err
	}
	if authChallenge.Error != "" {
		return "", errors.New(authChallenge.Error)
	}

	res, err := a.svc.Auth(context.Background(), &admin.Auth_Request{
		ChallengeId: authChallenge.ChallengeId,
		Signature:   ed25519.Sign(privateKey, authChallenge.Challenge),
	})
	if err != nil {
		return "", err
	}
	if res.Error != "" {
		return "", errors.New(res.Error)
	}
	a.sessionToken = res.SessionToken
	return res.SessionToken, nil
}

func (a *AdminClient) UpdateUserFingerprint(username, fingerprint string) error {
	res, err := a.svc.UpdateUserFingerprint(context.Background(), &admin.UpdateUserFingerprint_Request{
		SessionToken: a.sessionToken,
		Username:     username,
		Fingerprint:  fingerprint,
	})
	if err != nil {
		return err
	}
	if res.Error != "" {
		return errors.New(res.Error)
	}
	return nil
}
func (a *AdminClient) SendMessageToRoom(room, message string) error {
	res, err := a.svc.SendMessageToRoom(context.Background(), &admin.SendMessageToRoom_Request{
		SessionToken: a.sessionToken,
		Room:         room,
		Message:      message,
	})
	if err != nil {
		return err
	}
	if res.Error != "" {
		return errors.New(res.Error)
	}
	return nil
}

func main() {
	log.Println("reading private key...")
	data, err := os.ReadFile("./client-key")
	if err != nil {
		log.Fatal(err)
	}
	signer, err := ssh.ParsePrivateKey(data)
	if err != nil {
		log.Fatal(err)
	}

	//_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	//if err != nil {
	//	log.Fatal(err)
	//}

	client, err := ssh.Dial("tcp", "127.0.0.1:2222", &ssh.ClientConfig{
		Config:          ssh.Config{},
		User:            "client",
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		log.Fatal(err)
	}

	// RPC test
	channel, reqs, err := client.OpenChannel("rpc", nil)
	if err != nil {
		log.Println(err)
		return
	}
	go ssh.DiscardRequests(reqs)

	defer channel.Close()
	grpcConn, err := grpc.Dial("", grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
		return &sshnet.Conn{Channel: channel}, nil
	}), grpc.WithInsecure())
	if err != nil {
		log.Println(err)
		return
	}
	defer grpcConn.Close()
	adminClient := NewAdminClient(admin.NewAdminServiceClient(grpcConn))
	token, err := adminClient.Auth()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(token)
	log.Println(adminClient.UpdateUserFingerprint("chris", "SHA256:GLk/mTXZktyg18DbFzQdbl3dTFG4YHlO48HckkyJSt4"))
	log.Println(adminClient.SendMessageToRoom("default", "hello from the rpc interface"))

	// session stuff
	//session, err := client.NewSession()
	//if err != nil {
	//	log.Fatal(err)
	//}

	//out, err := session.StdoutPipe()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//in, err := session.StdinPipe()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//err = session.RequestPty("xterm", 40, 80, ssh.TerminalModes{
	//	ssh.ECHO:  0,
	//	ssh.IGNCR: 1,
	//})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//err = session.Shell()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Fprintf(in, ":clown:\n\r")
	//fmt.Fprintf(in, "/history chris\n\r")
	//done := make(chan struct{})
	//go func() {
	//	outScanner := bufio.NewScanner(out)
	//	for outScanner.Scan() {
	//		fmt.Println(outScanner.Text())
	//	}
	//	close(done)
	//}()
	//c := make(chan os.Signal, 1)
	//signal.Notify(c, os.Interrupt)
	//<-c
	//log.Println("exiting...")
	//session.Close()
	//<-done
}
