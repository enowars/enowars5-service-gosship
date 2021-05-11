package client

import (
	"context"
	"gosship/pkg/sshnet"
	"net"

	"golang.org/x/crypto/ssh"
	"google.golang.org/grpc"
)

func CreateNewGRPCClient(channel ssh.Channel) (*grpc.ClientConn, error) {
	grpcConn, err := grpc.Dial("", grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
		return &sshnet.Conn{Channel: channel}, nil
	}), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return grpcConn, nil
}
