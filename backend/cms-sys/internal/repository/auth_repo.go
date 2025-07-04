package repository

import (
	"errors"
	"github.com/google/uuid"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/types"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AuthRepository interface {
	CreateUser(user *types.CMSUser) error
	GetUserByEmail(email string) (*types.CMSUser, error)
	GetUserByID(id uuid.UUID) (*types.CMSUser, error)
	UpdateUser(user *types.CMSUser) error
	DeleteUser(id uuid.UUID) error
	EmailExists(email string) (bool, error)
}

type Repo struct {
	logger *logrus.Logger
	db     *gorm.DB
}

var _ AuthRepository = (*Repo)(nil)

func NewRepo(logger *logrus.Logger, db *gorm.DB) *Repo {
	return &Repo{
		logger: logger,
		db:     db,
	}
}

func (r *Repo) CreateUser(user *types.CMSUser) error {
	if err := r.db.Create(user).Error; err != nil {
		r.logger.WithError(err).Error("Failed to create user")
		return err
	}
	return nil
}

func (r *Repo) GetUserByEmail(email string) (*types.CMSUser, error) {
	var user types.CMSUser
	if err := r.db.Where("cms_user_email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		r.logger.WithError(err).Error("Failed to get user by email")
		return nil, err
	}
	return &user, nil
}

func (r *Repo) GetUserByID(id uuid.UUID) (*types.CMSUser, error) {
	var user types.CMSUser
	if err := r.db.Where("cms_user_id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		r.logger.WithError(err).Error("Failed to get user by ID")
		return nil, err
	}
	return &user, nil
}

func (r *Repo) UpdateUser(user *types.CMSUser) error {
	if err := r.db.Save(user).Error; err != nil {
		r.logger.WithError(err).Error("Failed to update user")
		return err
	}
	return nil
}

func (r *Repo) DeleteUser(id uuid.UUID) error {
	if err := r.db.Delete(&types.CMSUser{}, "cms_user_id = ?", id).Error; err != nil {
		r.logger.WithError(err).Error("Failed to delete user")
		return err
	}
	return nil
}

func (r *Repo) EmailExists(email string) (bool, error) {
	var count int64
	if err := r.db.Model(&types.CMSUser{}).Where("cms_user_email = ?", email).Count(&count).Error; err != nil {
		r.logger.WithError(err).Error("Failed to check if email exists")
		return false, err
	}
	return count > 0, nil
}
