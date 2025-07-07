package repository

import (
	"github.com/multi-tenants-cms-golang/cms-sys/internal/types"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OwnerRepository interface {
	CreateOwner(owner *types.CMSUser) error
	GetById(id string) (*types.CMSUser, error)
	UpdateOwner(owner *types.CMSUser) error
	GetOwnerByID(id string) (*types.CMSUser, error)
}

type OwnerRepositoryImpl struct {
	logger *logrus.Logger
	db     *gorm.DB
}

var _ OwnerRepository = (*OwnerRepositoryImpl)(nil)

func NewOwnerRepository(logger *logrus.Logger, db *gorm.DB) OwnerRepository {
	return &OwnerRepositoryImpl{logger: logger, db: db}
}

func (r *OwnerRepositoryImpl) CreateOwner(owner *types.CMSUser) error {
	if err := r.db.Create(owner).Error; err != nil {
		r.logger.WithError(err).Error("Failed to create owner")
		return err
	}
	return nil
}

func (r *OwnerRepositoryImpl) GetById(id string) (*types.CMSUser, error) {
	var owner types.CMSUser
	err := r.db.Where("cms_user_id", id).First(&owner).Error
	return &owner, err
}

func (r *OwnerRepositoryImpl) UpdateOwner(owner *types.CMSUser) error {
	if err := r.db.Save(owner).Error; err != nil {
		r.logger.WithError(err).Error("Failed to update owner")
	}
	return nil
}

func (r *OwnerRepositoryImpl) GetOwnerByID(id string) (*types.CMSUser, error) {
	var owner types.CMSUser
	err := r.db.Joins("JOIN cms_whole_sys_role r ON r.role_id = cms_user.cms_user_role").
		Where("cms_user_id = ? AND r.role_name = ?", id, string(types.CMSCustomer)).
		Find(&owner).Error
	return &owner, err
}
