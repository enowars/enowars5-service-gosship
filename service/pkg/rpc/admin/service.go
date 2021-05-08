package admin

import (
	"context"
	"gosship/pkg/database"

	"github.com/sirupsen/logrus"
)

type Service struct {
	log *logrus.Logger
	db  *database.Database
	UnimplementedAdminServiceServer
}

func NewService(log *logrus.Logger, db *database.Database) *Service {
	return &Service{
		log: log,
		db:  db,
	}
}

func (a *Service) Auth(ctx context.Context, request *Auth_Request) (*Auth_Response, error) {
	a.log.Println("auth-request", request.Config)
	return &Auth_Response{
		Error: "not implemented",
	}, nil
}
