package service

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type UserStore interface {
	UpdateEmailVerification(ctx context.Context, email string) error
	BeginTx(ctx context.Context) (pgx.Tx, error)
}

type AuthenticationService struct {
	store       UserStore
	logger      *logrus.Logger
	databaseCtx context.Context
}

func NewAuthenticationService(store UserStore, logger *logrus.Logger) *AuthenticationService {
	return &AuthenticationService{
		store:       store,
		logger:      logger,
		databaseCtx: context.Background(),
	}
}

//func (s *AuthenticationService) wholeRegistrationFlow(req *authenticationpb.RegisterRequest) (*authenticationpb.RegisterResponse, error) {
//	userToBeCreated := authConvert.ProtoToCreateParams(req)
//	err := s.store.xecTx(s.databaseCtx, func(tx pgx.Tx) error {
//		user, err := s.store.CreateUser(s.databaseCtx, *userToBeCreated)
//		if err != nil {
//			s.logger.WithFields(logrus.Fields{
//				"error": err.Error(),
//				"email": req.GetEmail(),
//			}).Error("failed to create user")
//			return err
//		}
//
//		ds, err := s.store.GetDefaultRoleIDs(s.databaseCtx)
//		if err != nil {
//			s.logger.WithFields(logrus.Fields{
//				"error": err.Error(),
//				"email": req.GetEmail(),
//			}).Error("failed to get default role IDs")
//			return err
//		}
//
//		s.store.AssignRolesToUser()
//		return nil
//	})
//	if err != nil {
//		return nil, err
//	}
//
//}
