package admin

import (
	"context"
	"errors"
	"gosship/pkg/chat"
	"gosship/pkg/database"
	"gosship/pkg/rpc/admin/auth"
	"sync"

	"github.com/sirupsen/logrus"
)

type Service struct {
	log              *logrus.Logger
	db               *database.Database
	host             *chat.Host
	sessionMu        sync.Mutex
	sessions         map[string]bool
	authChallengesMu sync.Mutex
	authChallenges   map[string][]byte
	UnimplementedAdminServiceServer
}

var ErrInvalidAuthChallenge = errors.New("invalid auth challenge")
var ErrInvalidAuthSignature = errors.New("invalid auth signature")
var ErrInvalidSessionToken = errors.New("invalid session token")

func NewService(log *logrus.Logger, db *database.Database, host *chat.Host) *Service {
	return &Service{
		log:            log,
		db:             db,
		host:           host,
		sessions:       make(map[string]bool),
		authChallenges: make(map[string][]byte),
	}
}

func (s *Service) GetAuthChallenge(ctx context.Context, request *GetAuthChallenge_Request) (*GetAuthChallenge_Response, error) {
	id, challenge := auth.CreateAuthChallenge()
	s.authChallengesMu.Lock()
	defer s.authChallengesMu.Unlock()
	s.authChallenges[id] = challenge
	return &GetAuthChallenge_Response{
		ChallengeId: id,
		Challenge:   challenge,
	}, nil
}

func (s *Service) Auth(ctx context.Context, request *Auth_Request) (*Auth_Response, error) {
	s.authChallengesMu.Lock()
	defer s.authChallengesMu.Unlock()
	challenge, hasChallenge := s.authChallenges[request.ChallengeId]
	if !hasChallenge {
		return nil, ErrInvalidAuthChallenge
	}
	if !auth.VerifySignature(challenge, request.Signature) {
		return nil, ErrInvalidAuthSignature
	}
	delete(s.authChallenges, request.ChallengeId)
	sessionToken := auth.GenerateRandomSessionToken()
	s.sessionMu.Lock()
	defer s.sessionMu.Unlock()
	s.sessions[sessionToken] = true
	return &Auth_Response{
		SessionToken: sessionToken,
	}, nil
}

func (s *Service) checkSession(token string) error {
	s.sessionMu.Lock()
	defer s.sessionMu.Unlock()
	if !s.sessions[token] {
		return ErrInvalidSessionToken
	}
	return nil
}

func (s *Service) UpdateUserFingerprint(ctx context.Context, request *UpdateUserFingerprint_Request) (*UpdateUserFingerprint_Response, error) {
	if err := s.checkSession(request.SessionToken); err != nil {
		return nil, err
	}
	err := s.db.UpdateUserFingerprint(request.Username, request.Fingerprint)
	if err != nil {
		return nil, err
	}
	return &UpdateUserFingerprint_Response{}, nil
}

func (s *Service) SendMessageToRoom(ctx context.Context, request *SendMessageToRoom_Request) (*SendMessageToRoom_Response, error) {
	if err := s.checkSession(request.SessionToken); err != nil {
		return nil, nil
	}
	s.host.RoomAnnouncement(request.Room, request.Message)
	return &SendMessageToRoom_Response{}, nil
}

func (s *Service) DumpMessages(request *DumpMessages_Request, server AdminService_DumpMessagesServer) error {
	//server.Send()
	return nil
}
