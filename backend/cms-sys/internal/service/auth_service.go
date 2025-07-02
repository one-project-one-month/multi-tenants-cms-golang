package service

import (
	"github.com/multi-tenants-cms-golang/cms-sys/internal/repository"
	"github.com/sirupsen/logrus"
)

type AuthService interface {
	Login(email, password string) (string, error)
	Register(email, password string) error
}

type Service struct {
	log  *logrus.Logger
	repo repository.AuthRepository
}

var _ (AuthService) = (*Service)(nil)

func NewService(
	log *logrus.Logger,
	repo repository.AuthRepository,
) *Service {
	return &Service{
		log:  log,
		repo: repo,
	}
}

func (s Service) Login(email, password string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) Register(email, password string) error {
	//TODO implement me
	panic("implement me")
}
