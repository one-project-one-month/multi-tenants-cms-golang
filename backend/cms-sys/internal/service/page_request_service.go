package service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
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
	ChangeStatus(req types.ChangeStatusPageRequest, currentUserID uuid.UUID) error
	ApprovePageRequest(req *types.ApprovePageRequest) (err error)
}

type PageRequestServiceImpl struct {
	logger      *logrus.Logger
	repo        repository.PageRequestRepository
	pageService PageService
	ownerRepo   repository.OwnerRepository
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

func (s *PageRequestServiceImpl) ChangeStatus(req types.ChangeStatusPageRequest, currentUserID uuid.UUID) error {
	requestUUID, err := uuid.Parse(req.RequestID)
	if err != nil {
		return errors.New("invalid request ID format")
	}

	_, err = s.repo.GetById(requestUUID)
	if err != nil || errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("invalid pageRequestId: %w", err)
	}

	if err := s.repo.UpdateStatus(requestUUID, req.Status); err != nil {
		return err
	}

	if req.Status == "APPROVED" {
		requestUUID, err := uuid.Parse(req.RequestID)
		if err != nil {
			return errors.New("invalid request ID format")
		}

		pageRequest, err := s.repo.GetById(requestUUID)
		if err != nil {
			return fmt.Errorf("failed to fetch page request: %w", err)
		}

		createReq := &types.PageCreateRequest{
			PageRequestID: requestUUID,
			Title:         pageRequest.Title,
			Content:       utils.SafeString(pageRequest.Description),
			ImageUrl:      utils.SafeString(pageRequest.LogoUrl),
			OwnerId:       pageRequest.OwnerID,
			PublisherId:   currentUserID,
			Status:        utils.StringPtr("PUBLISHED"),
		}

		_, err = s.pageService.Create(createReq)
		if err != nil {
			return fmt.Errorf("failed to auto-create page: %w", err)
		}

		owner, err := s.ownerRepo.GetById(pageRequest.OwnerID.String())
		if err != nil {
			return fmt.Errorf("failed to fetch owner: %w", err)
		}
		log.Printf("Owner fetched: %+v\n", owner) // TODO: Delete this log after you finish service to service communication

		// TODO: Ko Swan continue Communication with LMS Service here to create Tenant and send email to Owner
	}

	return nil
}
func (s *PageRequestServiceImpl) ApprovePageRequest(req *types.ApprovePageRequest) (err error) {
	//pageRequest , err := s.repo.GetById(req.PageRequestID)
	//if err != nil {
	//	if errors.Is(err, gorm.ErrRecordNotFound) {
	//		return errors.New("page request not found")
	//	}
	//}
	//s.repo
	//resp := utils.CreateTenant(req)

	return errors.New("invalid request ID format")
}
