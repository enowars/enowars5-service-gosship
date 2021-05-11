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

func NewGRPCServer(log *logrus.Logger, db *database.Database, host *chat.Host) *GRPCServer {
	grpcServer := grpc.NewServer()
	adminService := admin.NewService(log, db, host)
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
	s.adminService.PrepareForSession(ctx.SessionID())
	s.listener.PushChannel(ch)
}

func (s *GRPCServer) Serve() {
	err := s.grpcServer.Serve(s.listener)
	if err != nil {
		s.log.Error(err)
	}
}
