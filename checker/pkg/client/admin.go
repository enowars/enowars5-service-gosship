package client

import (
	"checker/service/admin"
	"checker/service/database"
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"io"

	"golang.org/x/crypto/ssh"
	"google.golang.org/grpc"
)

var PrivateKey ed25519.PrivateKey

func init() {
	privateKeyRaw, err := hex.DecodeString("215c8787c1b079be149db3da5e297a9b39ff008dee69b2b9f115d51d4547664580de0e58c0842f83cf95f9772a5a13c167dd4c0e3fd02913076d16df828fbbb2")
	if err != nil {
		panic(err)
	}
	PrivateKey = privateKeyRaw
}

var ErrInvalidAuthChallenge = errors.New("invalid auth challenge")

type AdminClient struct {
	svc          admin.AdminServiceClient
	SessionToken string
	PublicKey    ssh.PublicKey
}

func NewAdminClient(grpcConn *grpc.ClientConn, publicKey ssh.PublicKey) *AdminClient {
	return &AdminClient{
		svc:       admin.NewAdminServiceClient(grpcConn),
		PublicKey: publicKey,
	}
}

func (a *AdminClient) Auth() (string, error) {
	authChallenge, err := a.svc.GetAuthChallenge(context.Background(), &admin.GetAuthChallenge_Request{})
	if err != nil {
		return "", err
	}

	if len(authChallenge.Challenge) != 576 {
		return "", ErrInvalidAuthChallenge
	}

	err = a.PublicKey.Verify(authChallenge.Challenge[:512], &ssh.Signature{
		Format: "ssh-ed25519",
		Blob:   authChallenge.Challenge[512:],
	})
	if err != nil {
		return "", err
	}

	res, err := a.svc.Auth(context.Background(), &admin.Auth_Request{
		ChallengeId: authChallenge.ChallengeId,
		Signature:   ed25519.Sign(PrivateKey, authChallenge.Challenge),
	})
	if err != nil {
		return "", err
	}
	a.SessionToken = res.SessionToken
	return res.SessionToken, nil
}

func (a *AdminClient) SendMessageToRoom(room, message string) error {
	_, err := a.svc.SendMessageToRoom(context.Background(), &admin.SendMessageToRoom_Request{
		SessionToken: a.SessionToken,
		Room:         room,
		Message:      message,
	})
	if err != nil {
		return err
	}
	return nil
}

func (a *AdminClient) DumpDirectMessages(username string, emit func(entry *database.MessageEntry)) error {
	stream, err := a.svc.DumpDirectMessages(context.Background(), &admin.DumpDirectMessages_Request{
		SessionToken: a.SessionToken,
		Username:     username,
	})
	if err != nil {
		return err
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		emit(res.Message)
	}
	return nil
}
