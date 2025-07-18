package service

import (
	"errors"
	"fmt"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/mapper"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/repository"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/types"
	"github.com/multi-tenants-cms-golang/cms-sys/pkg/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"math"
)

type PageService interface {
	GetAllPages(req *types.PaginateRequest) ([]*types.PageResponse, *utils.Pagination, error)
	Create(req *types.PageCreateRequest) (*types.PageResponse, error)
	GetById(id string) (*types.PageResponse, error)
	UpdateById(id string, req *types.PageUpdateRequest) (*types.PageResponse, error)
	DeleteById(id string) error
}

type PageServiceImpl struct {
	log             *logrus.Logger
	repo            repository.PageRepository
	pageRequestRepo repository.PageRequestRepository
	ownerRepo       repository.OwnerRepository
}

var _ PageService = (*PageServiceImpl)(nil)

func NewPageService(
	log *logrus.Logger,
	repo repository.PageRepository,
	pageRequestRepo repository.PageRequestRepository,
	ownerRepo repository.OwnerRepository,
) PageService {
	return &PageServiceImpl{
		log:             log,
		repo:            repo,
		pageRequestRepo: pageRequestRepo,
		ownerRepo:       ownerRepo,
	}
}

/*
	func (p *PageServiceImpl) GetAllPages(req *types.PaginateRequest) ([]*types.PageResponse, *utils.Pagination, error) {
		total, err := p.repo.Count()
		if err != nil {
			p.log.WithError(err).Errorln("Failed to count number of pages")
			return nil, nil, err
		}

		offset := utils.CalculateOffset(req.Page, req.Limit)

		pages, err := p.repo.GetAllPages(offset, req.Limit)
		if err != nil {
			p.log.WithError(err).Errorln("Failed to get all pages")
			return nil, nil, err
		}

		var responses []*types.PageResponse
		for _, page := range pages {
			responses = append(responses, mapper.ToPageResponse(page))
		}

		pagination := &utils.Pagination{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      total,
			TotalPages: int(math.Ceil(float64(total) / float64(req.Limit))),
		}

		return responses, pagination, nil
	}
*/
func (p *PageServiceImpl) GetAllPages(req *types.PaginateRequest) ([]*types.PageResponse, *utils.Pagination, error) {
	total, err := p.repo.Count()
	if err != nil {
		p.log.WithError(err).Errorln("Failed to count number of pages")
		return nil, nil, err
	}

	offset := utils.CalculateOffset(req.Page, req.Limit)

	pages, err := p.repo.GetAllPagesWithRelations(offset, req.Limit)
	if err != nil {
		p.log.WithError(err).Errorln("Failed to get all pages with relations")
		return nil, nil, err
	}

	var responses []*types.PageResponse
	for _, page := range pages {
		responses = append(responses, mapper.ToPageResponse(page))
	}

	pagination := &utils.Pagination{
		Page:       req.Page,
		Limit:      req.Limit,
		Total:      total,
		TotalPages: int(math.Ceil(float64(total) / float64(req.Limit))),
	}

	return responses, pagination, nil
}

func (p *PageServiceImpl) Create(req *types.PageCreateRequest) (*types.PageResponse, error) {
	// Check provided page request id exists
	_, err := p.pageRequestRepo.GetById(req.PageRequestID)
	if err != nil || errors.Is(err, gorm.ErrRecordNotFound) {
		p.log.WithField("pageRequestId", req.PageRequestID).Warn("Page request not found")
		return nil, fmt.Errorf("invalid pageRequestId: %w", err)
	}
	// Check provided owner id exists
	_, err = p.ownerRepo.GetById(req.OwnerId.String())
	if err != nil || errors.Is(err, gorm.ErrRecordNotFound) {
		p.log.WithField("ownerId", req.OwnerId).Warn("Owner not found")
		return nil, fmt.Errorf("invalid ownerId: %w", err)
	}
	// Check provided publisher id exists
	_, err = p.ownerRepo.GetById(req.PublisherId.String())
	if err != nil || errors.Is(err, gorm.ErrRecordNotFound) {
		p.log.WithField("publisherId", req.PublisherId).Warn("Publisher not found")
		return nil, fmt.Errorf("invalid publisherId: %w", err)
	}

	status := types.PageStatus("DRAFT")
	if req.Status != nil {
		status = types.PageStatus(*req.Status)
	}

	page := &types.Page{
		PageRequestID:      req.PageRequestID,
		Title:              req.Title,
		Content:            req.Content,
		ImageURL:           &req.ImageUrl,
		Status:             status,
		OwnerID:            req.OwnerId,
		PublishedByStaffID: &req.PublisherId,
	}

	if err := p.repo.CreatePage(page); err != nil {
		p.log.WithError(err).Errorln("Failed to create page")
		return nil, errors.New("failed to create page")
	}

	return mapper.ToPageResponse(page), nil
}

func (p *PageServiceImpl) GetById(id string) (*types.PageResponse, error) {
	page, err := p.repo.GetById(id)

	if err != nil {
		p.log.WithError(err).Errorln("Failed to get page")
		return nil, errors.New("failed to get page")
	}

	return mapper.ToPageResponse(page), nil
}

func (p *PageServiceImpl) UpdateById(id string, req *types.PageUpdateRequest) (*types.PageResponse, error) {
	page, err := p.repo.GetById(id)
	if err != nil {
		p.log.WithError(err).Error("Failed to get page")
		return nil, errors.New("failed to get page")
	}

	// Check provided page request id exists
	_, err = p.pageRequestRepo.GetById(req.PageRequestID)
	if err != nil || errors.Is(err, gorm.ErrRecordNotFound) {
		p.log.WithField("pageRequestId", req.PageRequestID).Warn("Page request not found")
		return nil, fmt.Errorf("invalid pageRequestId: %w", err)
	}
	// Check provided owner id exists
	_, err = p.ownerRepo.GetById(req.OwnerId.String())
	if err != nil || errors.Is(err, gorm.ErrRecordNotFound) {
		p.log.WithField("ownerId", req.OwnerId).Warn("Owner not found")
		return nil, fmt.Errorf("invalid ownerId: %w", err)
	}
	// Check provided publisher id exists
	_, err = p.ownerRepo.GetById(req.PublisherId.String())
	if err != nil || errors.Is(err, gorm.ErrRecordNotFound) {
		p.log.WithField("publisherId", req.PublisherId).Warn("Publisher not found")
		return nil, fmt.Errorf("invalid publisherId: %w", err)
	}

	page.PageRequestID = req.PageRequestID
	page.Title = req.Title
	page.Content = req.Content
	page.ImageURL = &req.ImageUrl
	page.Status = types.PageStatus(req.Status)
	page.OwnerID = req.OwnerId
	page.PublishedByStaffID = &req.PublisherId

	if err := p.repo.UpdatePage(page); err != nil {
		p.log.WithError(err).Error("Failed to update page")
		return nil, errors.New("failed to update page")
	}

	return mapper.ToPageResponse(page), nil
}

func (p *PageServiceImpl) DeleteById(id string) error {
	_, err := p.repo.GetById(id)
	if err != nil {
		p.log.WithError(err).Errorln("Failed to get page")
		return errors.New("failed to get page")
	}
	err = p.repo.DeleteById(id)
	if err != nil {
		return errors.New("failed to delete page")
	}
	return nil
}
