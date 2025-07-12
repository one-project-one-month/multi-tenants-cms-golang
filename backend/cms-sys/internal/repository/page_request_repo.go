package repository

import (
	"github.com/google/uuid"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/types"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PageRequestRepository interface {
	CreatePageRequest(pageRequest *types.PageRequest) error
	GetAllPageRequests(offset, limit int) ([]*types.PageRequest, error)
	CountPageRequests() (int64, error)
	GetById(id uuid.UUID) (*types.PageRequest, error)
	UpdateStatus(id uuid.UUID, status types.RequestStatus) error
}

type PageRequestRepositoryImpl struct {
	logger *logrus.Logger
	db     *gorm.DB
}

var _ PageRequestRepository = (*PageRequestRepositoryImpl)(nil)

func NewPageRequestRepository(logger *logrus.Logger, db *gorm.DB) PageRequestRepository {
	return &PageRequestRepositoryImpl{
		logger: logger,
		db:     db,
	}
}

func (r *PageRequestRepositoryImpl) CreatePageRequest(pageRequest *types.PageRequest) error {
	if err := r.db.Create(pageRequest).Error; err != nil {
		r.logger.WithError(err).Error("Failed to create page request")
		return err
	}
	return nil
}

func (r *PageRequestRepositoryImpl) GetAllPageRequests(offset, limit int) ([]*types.PageRequest, error) {
	var pageRequests []*types.PageRequest
	if err := r.db.Offset(offset).Limit(limit).Find(&pageRequests).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get all page requests")
		return nil, err
	}
	return pageRequests, nil
}

func (r *PageRequestRepositoryImpl) CountPageRequests() (int64, error) {
	var count int64
	if err := r.db.Model(&types.PageRequest{}).Count(&count).Error; err != nil {
		r.logger.WithError(err).Error("Failed to count page requests")
		return 0, err
	}
	return count, nil
}


func (r *PageRequestRepositoryImpl) GetById(id uuid.UUID) (*types.PageRequest, error) {
	var pageRequest types.PageRequest
	if err := r.db.Where("request_id = ?", id).First(&pageRequest).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get page request")
		return nil, err
	}
	return &pageRequest, nil
}

func (r *PageRequestRepositoryImpl) UpdateStatus(id uuid.UUID, status types.RequestStatus) error {
    result := r.db.Model(&types.PageRequest{}).Where("request_id = ?", id).Update("status", status)
    if result.Error != nil {
        r.logger.WithError(result.Error).Error("Failed to update status of page request")
        return result.Error
    }
    if result.RowsAffected == 0 {
        return gorm.ErrRecordNotFound
    }
    return nil
}
