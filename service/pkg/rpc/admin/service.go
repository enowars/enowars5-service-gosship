package admin

import (
	"context"
	"errors"
	"gosship/pkg/chat"
	"gosship/pkg/database"
	"gosship/pkg/rpc/admin/auth"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/peer"
)

type AuthChallenge struct {
	Challenge []byte
	Session   string
	Timestamp time.Time
}

type Service struct {
	log              *logrus.Logger
	db               *database.Database
	host             *chat.Host
	sessionMu        sync.Mutex
	sessions         map[string]bool
	authChallengesMu sync.Mutex
	authChallenges   map[string]*AuthChallenge
	UnimplementedAdminServiceServer
}

var ErrInvalidRequest = errors.New("invalid request")
var ErrInvalidAuthChallenge = errors.New("invalid auth challenge")
var ErrExpiredAuthChallenge = errors.New("expired auth challenge")
var ErrInvalidSSHSession = errors.New("invalid ssh session")
var ErrInvalidAuthSignature = errors.New("invalid auth signature")
var ErrInvalidSessionToken = errors.New("invalid session token")

func NewService(log *logrus.Logger, db *database.Database, host *chat.Host) *Service {
	return &Service{
		log:            log,
		db:             db,
		host:           host,
		sessions:       make(map[string]bool),
		authChallenges: make(map[string]*AuthChallenge),
	}
}

func (s *Service) GetAuthChallenge(ctx context.Context, request *GetAuthChallenge_Request) (*GetAuthChallenge_Response, error) {
	peer, ok := peer.FromContext(ctx)
	if !ok {
		return nil, ErrInvalidRequest
	}
	id, challengePayload := auth.CreateAuthChallenge()
	challenge := &AuthChallenge{
		Challenge: challengePayload,
		Session:   peer.Addr.String(),
		Timestamp: time.Now(),
	}
	s.authChallengesMu.Lock()
	defer s.authChallengesMu.Unlock()
	s.authChallenges[id] = challenge
	s.log.Println("[admin-service] new auth challenge")
	return &GetAuthChallenge_Response{
		ChallengeId: id,
		Challenge:   challengePayload,
	}, nil
}

func (s *Service) Auth(ctx context.Context, request *Auth_Request) (*Auth_Response, error) {
	peer, ok := peer.FromContext(ctx)
	if !ok {
		return nil, ErrInvalidRequest
	}
	s.authChallengesMu.Lock()
	defer s.authChallengesMu.Unlock()
	challenge, hasChallenge := s.authChallenges[request.ChallengeId]
	if !hasChallenge {
		return nil, ErrInvalidAuthChallenge
	}
	if challenge.Session != peer.Addr.String() {
		return nil, ErrInvalidSSHSession
	}

	if challenge.Timestamp.Before(time.Now().Add(-10 * time.Second)) {
		delete(s.authChallenges, request.ChallengeId)
		return nil, ErrExpiredAuthChallenge
	}
	if !auth.VerifySignature(challenge.Challenge, request.Signature) {
		return nil, ErrInvalidAuthSignature
	}

	delete(s.authChallenges, request.ChallengeId)
	sessionToken := auth.GenerateRandomSessionToken()
	s.sessionMu.Lock()
	defer s.sessionMu.Unlock()
	s.sessions[sessionToken] = true
	s.log.Println("[admin-service] new session")
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

func (s *Service) SendMessageToRoom(ctx context.Context, request *SendMessageToRoom_Request) (*SendMessageToRoom_Response, error) {
	if err := s.checkSession(request.SessionToken); err != nil {
		return nil, err
	}
	s.log.Printf("[admin-service] new message for room %s: %s", request.Room, request.Message)
	s.host.RoomAnnouncement(request.Room, request.Message)
	return &SendMessageToRoom_Response{}, nil
}

func (s *Service) DumpDirectMessages(request *DumpDirectMessages_Request, server AdminService_DumpDirectMessagesServer) error {
	if err := s.checkSession(request.SessionToken); err != nil {
		return err
	}
	s.log.Printf("[admin-service] direct messages requested for %s", request.Username)
	return s.db.DumpDirectMessages(request.Username, func(entry *database.MessageEntry) error {
		return server.Send(&DumpDirectMessages_Response{Message: entry})
	})
}
