package client

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"gosship/pkg/rpc/admin"

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

type AdminClient struct {
	svc          admin.AdminServiceClient
	SessionToken string
}

func NewAdminClient(grpcConn *grpc.ClientConn) *AdminClient {
	return &AdminClient{svc: admin.NewAdminServiceClient(grpcConn)}
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
		Signature:   ed25519.Sign(PrivateKey, authChallenge.Challenge),
	})
	if err != nil {
		return "", err
	}
	if res.Error != "" {
		return "", errors.New(res.Error)
	}
	a.SessionToken = res.SessionToken
	return res.SessionToken, nil
}

func (a *AdminClient) UpdateUserFingerprint(username, fingerprint string) error {
	res, err := a.svc.UpdateUserFingerprint(context.Background(), &admin.UpdateUserFingerprint_Request{
		SessionToken: a.SessionToken,
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
		SessionToken: a.SessionToken,
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
