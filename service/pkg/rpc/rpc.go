package rpc

import (
	"gosship/pkg/chat"
	"gosship/pkg/database"
	"gosship/pkg/rpc/admin"
	"gosship/pkg/sshnet"

	"github.com/gliderlabs/ssh"
	"github.com/sirupsen/logrus"
	gossh "golang.org/x/crypto/ssh"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	log          *logrus.Logger
	db           *database.Database
	host         *chat.Host
	listener     *sshnet.Listener
	grpcServer   *grpc.Server
	adminService *admin.Service
}

func NewGRPCServer(log *logrus.Logger, db *database.Database, host *chat.Host, signer gossh.Signer) *GRPCServer {
	grpcServer := grpc.NewServer()
	adminService := admin.NewService(log, db, host, signer)
	admin.RegisterAdminServiceServer(grpcServer, adminService)
	return &GRPCServer{
		log:          log,
		db:           db,
		host:         host,
		listener:     sshnet.NewListener(),
		grpcServer:   grpcServer,
		adminService: adminService,
	}
}

func (s *GRPCServer) Handle(srv *ssh.Server, conn *gossh.ServerConn, newChan gossh.NewChannel, ctx ssh.Context) {
	ch, reqs, err := newChan.Accept()
	if err != nil {
		s.log.Error(err)
		return
	}
	go gossh.DiscardRequests(reqs)
	s.log.Printf("new rpc channel opened by %s (%s)", ctx.User(), ctx.RemoteAddr())
	s.listener.PushChannel(ch, ctx.SessionID())
}

func (s *GRPCServer) Serve() {
	err := s.grpcServer.Serve(s.listener)
	if err != nil {
		s.log.Error(err)
	}
}
