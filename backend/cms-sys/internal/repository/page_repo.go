package repository

import (
	"github.com/multi-tenants-cms-golang/cms-sys/internal/types"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PageRepository interface {
	GetAllPages(offset, limit int) ([]*types.Page, error)
	Count() (int64, error)
	CreatePage(page *types.Page) error
	UpdatePage(page *types.Page) error
	GetById(id string) (*types.Page, error)
	DeleteById(id string) error
}

type PageRepositoryImpl struct {
	logger *logrus.Logger
	db     *gorm.DB
}

var _ PageRepository = (*PageRepositoryImpl)(nil)

func NewPageRepository(log *logrus.Logger, db *gorm.DB) PageRepository {
	return &PageRepositoryImpl{
		logger: log,
		db:     db,
	}
}

func (p *PageRepositoryImpl) GetAllPages(offset, limit int) ([]*types.Page, error) {
	var pages []*types.Page
	if err := p.db.Offset(offset).Limit(limit).Find(&pages).Error; err != nil {
		p.logger.WithError(err).Error("failed to fetch pages")
		return nil, err
	}
	return pages, nil
}

func (p *PageRepositoryImpl) Count() (int64, error) {
	var count int64
	if err := p.db.Model(&types.Page{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (p *PageRepositoryImpl) CreatePage(page *types.Page) error {
	if err := p.db.Create(page).Error; err != nil {
		p.logger.WithError(err).Error("failed to create page")
		return err
	}
	return nil
}

func (p *PageRepositoryImpl) UpdatePage(page *types.Page) error {
	if err := p.db.Save(page).Error; err != nil {
		p.logger.WithError(err).Error("failed to update page")
		return err
	}
	return nil
}

func (p *PageRepositoryImpl) GetById(id string) (*types.Page, error) {
	var page types.Page
	if err := p.db.Where("page_id", id).First(&page).Error; err != nil {
		p.logger.WithError(err).Error("failed to fetch page")
		return nil, err
	}
	return &page, nil
}

func (p *PageRepositoryImpl) DeleteById(id string) error {
	if err := p.db.Delete(&types.Page{}, "page_id = ?", id).Error; err != nil {
		p.logger.WithError(err).Error("failed to delete page")
		return err
	}
	return nil
}
