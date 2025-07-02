package repository

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AuthRepository interface {
	Login(email string) (string, error)
	Register(email, password string) error
}

type Repo struct {
	logger *logrus.Logger
	db     *gorm.DB
}

func (r Repo) Login(email string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (r Repo) Register(email, password string) error {
	//TODO implement me
	panic("implement me")
}

var _ (AuthRepository) = (*Repo)(nil)
