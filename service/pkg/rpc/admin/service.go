package admin

import (
	"context"
	"gosship/pkg/database"
	"gosship/pkg/rpc/admin/auth"
	"sync"

	"github.com/sirupsen/logrus"
)

type Service struct {
	log       *logrus.Logger
	db        *database.Database
	sessionMu sync.Mutex
	sessions  map[string]bool
	UnimplementedAdminServiceServer
}

func NewService(log *logrus.Logger, db *database.Database) *Service {
	return &Service{
		log:      log,
		db:       db,
		sessions: make(map[string]bool),
	}
}

func (s *Service) Auth(ctx context.Context, request *Auth_Request) (*Auth_Response, error) {
	if !auth.CheckPassword(request.Password) {
		return &Auth_Response{Error: "invalid password"}, nil
	}
	sessionToken := auth.GenerateRandomSessionToken()
	s.sessionMu.Lock()
	defer s.sessionMu.Unlock()
	s.sessions[sessionToken] = true
	return &Auth_Response{
		SessionToken: sessionToken,
	}, nil
}

func (s *Service) ResetUserFingerprint(ctx context.Context, request *ResetUserFingerprint_Request) (*ResetUserFingerprint_Response, error) {
	s.sessionMu.Lock()
	defer s.sessionMu.Unlock()
	if !s.sessions[request.SessionToken] {
		return &ResetUserFingerprint_Response{Error: "invalid session token"}, nil
	}
	//TODO: s.db.ResetUserFingerprint
	return &ResetUserFingerprint_Response{Error: "TODO"}, nil
}
