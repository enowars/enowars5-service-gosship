package rpc

import (
	"gosship/pkg/database"
	"gosship/pkg/rpc/admin"
	"gosship/pkg/sshnet"

	"github.com/gliderlabs/ssh"
	"github.com/sirupsen/logrus"
	gossh "golang.org/x/crypto/ssh"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	log        *logrus.Logger
	db         *database.Database
	listener   *sshnet.Listener
	grpcServer *grpc.Server
}

func NewGRPCServer(log *logrus.Logger, db *database.Database) *GRPCServer {
	grpcServer := grpc.NewServer()
	admin.RegisterAdminServiceServer(grpcServer, admin.NewService(log, db))
	return &GRPCServer{
		log:        log,
		db:         db,
		listener:   sshnet.NewListener(),
		grpcServer: grpcServer,
	}
}

func (s *GRPCServer) Handle(srv *ssh.Server, conn *gossh.ServerConn, newChan gossh.NewChannel, ctx ssh.Context) {
	ch, reqs, err := newChan.Accept()
	if err != nil {
		s.log.Error(err)
		return
	}
	go gossh.DiscardRequests(reqs)
	s.listener.PushChannel(ch)
}

func (s *GRPCServer) Serve() {
	err := s.grpcServer.Serve(s.listener)
	if err != nil {
		s.log.Error(err)
	}
}
