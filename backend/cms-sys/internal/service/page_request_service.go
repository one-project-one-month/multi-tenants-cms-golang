package service

import (
	"errors"
	"github.com/google/uuid"
	"math"
	"strings"

	"github.com/multi-tenants-cms-golang/cms-sys/internal/mapper"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/repository"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/types"
	"github.com/multi-tenants-cms-golang/cms-sys/pkg/utils"
	"github.com/sirupsen/logrus"
	"time"
)

type PageRequestService interface {
	CreatePageRequest(req types.CreatePageRequest, logourl *string) (*types.PageRequestResponse, error)
	GetAllPageRequests(req *types.PaginateRequest) ([]*types.PageRequestResponse, *utils.Pagination, error)
}

type PageRequestServiceImpl struct {
	logger *logrus.Logger
	repo   repository.PageRequestRepository
}

var _ PageRequestService = (*PageRequestServiceImpl)(nil)

func NewPageRequestService(logger *logrus.Logger, repo repository.PageRequestRepository) *PageRequestServiceImpl {
	return &PageRequestServiceImpl{
		logger: logger,
		repo:   repo,
	}
}

func (s *PageRequestServiceImpl) CreatePageRequest(req types.CreatePageRequest, logoUrl *string) (*types.PageRequestResponse, error) {
	ownerUUID, err := uuid.Parse(strings.TrimSpace(req.OwnerID))
	if err != nil {
		return nil, errors.New("invalid owner ID format")
	}

	pageRequest := &types.PageRequest{
		RequestID:   uuid.New(),
		OwnerID:     ownerUUID,
		RequestType: req.RequestType,
		Title:       req.Title,
		Description: &req.Description,
		PageUrl:     req.PageUrl,
		LogoUrl:     logoUrl,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.CreatePageRequest(pageRequest); err != nil {
		s.logger.WithError(err).Error("Failed to create page request")
		return nil, err
	}

	return mapper.ToPageRequestResponse(pageRequest), nil
}

func (s *PageRequestServiceImpl) GetAllPageRequests(req *types.PaginateRequest) ([]*types.PageRequestResponse, *utils.Pagination, error) {
	total, err := s.repo.CountPageRequests()
	if err != nil {
		s.logger.WithError(err).Error("Failed to count page requests")
		return nil, nil, err
	}

	offset := utils.CalculateOffset(req.Page, req.Limit)

	pageRequests, err := s.repo.GetAllPageRequests(offset, req.Limit)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get all page requests")
		return nil, nil, err
	}

	var responses []*types.PageRequestResponse
	for _, pr := range pageRequests {
		responses = append(responses, mapper.ToPageRequestResponse(pr))
	}

	pagination := &utils.Pagination{
		Page:       req.Page,
		Limit:      req.Limit,
		Total:      total,
		TotalPages: int(math.Ceil(float64(total) / float64(req.Limit))),
	}

	return responses, pagination, nil
}
