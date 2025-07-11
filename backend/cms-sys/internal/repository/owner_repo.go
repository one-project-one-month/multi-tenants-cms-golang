package repository

import (
	"github.com/multi-tenants-cms-golang/cms-sys/internal/types"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OwnerRepository interface {
	GetAllOwners() ([]types.CMSUser, error)
	CreateOwner(owner *types.CMSUser) error
	GetById(id string) (*types.CMSUser, error)
	UpdateOwner(owner *types.CMSUser) error
	GetOwnerByID(id string) (*types.CMSUser, error)
	DeleteOwnerByID(id string) error
	OwnerHasAssociations(id string) bool
	ForceDeleteOwner(id string) error
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

func (r *OwnerRepositoryImpl) GetAllOwners() ([]types.CMSUser, error) {
	var owners []types.CMSUser
	if err := r.db.Find(&owners).Error; err != nil {
		return nil, err
	}
	return owners, nil
}

func (r *OwnerRepositoryImpl) GetOwnerByID(id string) (*types.CMSUser, error) {
	var owner types.CMSUser
	err := r.db.Joins("JOIN cms_whole_sys_role r ON r.role_id = cms_user.cms_user_role").
		Where("cms_user_id = ? AND r.role_name = ?", id, string(types.CMSCustomer)).
		Find(&owner).Error
	return &owner, err
}

func (r *OwnerRepositoryImpl) DeleteOwnerByID(id string) error {
	if err := r.db.Where("cms_user_id = ?", id).Delete(&types.CMSUser{}).Error; err != nil {
		r.logger.WithError(err).Error("Failed to delete owner by ID")
		return err
	}
	return nil
}

func (r *OwnerRepositoryImpl) OwnerHasAssociations(id string) bool {
	var count int64

	rawSQL := `
		SELECT 1 FROM cms_page WHERE owner_id = ? LIMIT 1
		UNION ALL
		SELECT 1 FROM cms_page_request WHERE owner_id = ? LIMIT 1
		UNION ALL
		SELECT 1 FROM cms_cus_purchase WHERE cms_cus_id = ? LIMIT 1
		UNION ALL
		SELECT 1 FROM user_page_request WHERE user_id = ? LIMIT 1
		LIMIT 1
	`
	if err := r.db.Raw(rawSQL, id, id, id, id).Scan(&count).Error; err != nil {
		r.logger.WithError(err).Error("Failed to check owner associations")
		return true
	}
	return count > 0
}

func (r *OwnerRepositoryImpl) ForceDeleteOwner(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("owner_id = ?", id).Delete(&types.Page{}).Error; err != nil {
			r.logger.WithError(err).Error("Failed to delete owned pages")
			return err
		}

		// Set Null staff publication
		if err := tx.Model(&types.Page{}).Where("published_by_staff_id = ?", id).
			Update("published_by_staff_id", nil).Error; err != nil {
			r.logger.WithError(err).Error("Failed to nullify published_by_staff_id")
			return err
		}

		// Delete Page Requests
		if err := tx.Where("owner_id = ?", id).Delete(&types.PageRequest{}).Error; err != nil {
			r.logger.WithError(err).Error("Failed to delete page requests")
			return err
		}

		if err := tx.Model(&types.PageRequest{}).Where("admin_id = ?", id).
			Update("admin_id", nil).Error; err != nil {
			r.logger.WithError(err).Error("Failed to nullify admin_id")
			return err
		}

		if err := tx.Where("cms_cus_id = ?", id).Delete(&types.CMSCusPurchase{}).Error; err != nil {
			r.logger.WithError(err).Error("Failed to delete purchases")
			return err
		}

		if err := tx.Where("user_id = ?", id).Delete(&types.UserPageRequest{}).Error; err != nil {
			r.logger.WithError(err).Error("Failed to delete user page requests")
			return err
		}

		if err := tx.Where("cms_user_id = ?", id).Delete(&types.CMSUser{}).Error; err != nil {
			r.logger.WithError(err).Error("Failed to delete CMS user")
			return err
		}

		return nil
	})
}

