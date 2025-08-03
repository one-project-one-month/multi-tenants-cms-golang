package service

import (
	"errors"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/mapper"
	"time"

	"github.com/google/uuid"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/repository"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/types"
	"github.com/multi-tenants-cms-golang/cms-sys/pkg/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OwnerService interface {
	Create(req types.OwnerCreateRequest) (*types.OwnerResponse, error)
	Update(id string, req types.OwnerUpdateRequest) (*types.OwnerResponse, error)
	GetAllOwners() ([]*types.OwnerResponse, error)
	GetOwnerByID(id string) (*types.OwnerResponse, error)
	BulkDeleteOwners(ids []string, force bool) error
}

type OwnerServiceImpl struct {
	log      *logrus.Logger
	repo     repository.OwnerRepository
	authRepo repository.AuthRepository
}

var _ OwnerService = (*OwnerServiceImpl)(nil)

func NewOwnerService(log *logrus.Logger, repo repository.OwnerRepository, authRepo repository.AuthRepository) OwnerService {
	return &OwnerServiceImpl{
		log:      log,
		repo:     repo,
		authRepo: authRepo,
	}
}

func (os *OwnerServiceImpl) Create(req types.OwnerCreateRequest) (*types.OwnerResponse, error) {
	exists, err := os.authRepo.EmailExists(req.Email)
	if err != nil {
		os.log.WithError(err).Error("Failed to check email exists!")
	}

	if exists {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		os.log.WithError(err).Error("Failed to hash password!")
		return nil, errors.New("filed to process password")
	}

	owner := &types.CMSUser{
		CMSUserID:    uuid.New(),
		CMSUserName:  req.Name,
		CMSUserEmail: req.Email,
		CMSNameSpace: &req.NameSpace, // TODO : Gotta fix it later as business logic
		Password:     hashedPassword,
		CMSUserRole:  string(types.CMSCustomer),
		Verified:     false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := os.repo.CreateOwner(owner); err != nil {
		os.log.WithError(err).Error("Failed to create owner")
		return nil, errors.New("failed to create owner")
	}

	createdOwner, err := os.repo.GetOwnerByID(owner.CMSUserID.String())
	if err != nil {
		os.log.WithError(err).Error("Failed to retrieve newly created owner")
		return nil, err
	}

	ownerResponse := mapper.ToOwnerResponse(createdOwner)
	return ownerResponse, nil
}

func (os *OwnerServiceImpl) Update(id string, req types.OwnerUpdateRequest) (*types.OwnerResponse, error) {
	owner, err := os.repo.GetById(id)
	if err != nil {
		os.log.WithError(err).Error("Failed to get owner with id ", id)
		return nil, errors.New("failed to get owner")
	}
	owner.CMSUserName = req.Name
	// TODO : Gotta fix it later as business logic
	owner.CMSNameSpace = &req.NameSpace
	if err := os.repo.UpdateOwner(owner); err != nil {
		os.log.WithError(err).Error("Failed to update owner")
		return nil, errors.New("failed to update owner")
	}

	ownerResponse := mapper.ToOwnerResponse(owner)
	return ownerResponse, nil
}

/*
func (os *OwnerServiceImpl) GetAllOwners() ([]types.OwnerResponse, error) {
	owners, err := os.repo.GetAllOwners()
	if err != nil {
		os.log.WithError(err).Error("Failed to get all owners")
		return nil, errors.New("failed to retrieve owners")
	}

	responses := make([]types.OwnerResponse, len(owners))

	for i, owner := range owners {
		var nameSpace string
		if owner.CMSNameSpace != nil {
			nameSpace = *owner.CMSNameSpace
		}

		responses[i] = types.OwnerResponse{
			ID:        owner.CMSUserID,
			Name:      owner.CMSUserName,
			Email:     owner.CMSUserEmail,
			Role:      owner.CMSUserRole,
			NameSpace: nameSpace,
			Verified:  owner.Verified,
		}
	}

	return responses, nil
}
*/

func (os *OwnerServiceImpl) GetAllOwners() ([]*types.OwnerResponse, error) {
	owners, err := os.repo.GetAllOwners()
	if err != nil {
		os.log.WithError(err).Error("Failed to get all owners from repository")
		return nil, err
	}

	return mapper.ToOwnerListResponse(owners), nil
}

func (os *OwnerServiceImpl) GetOwnerByID(id string) (*types.OwnerResponse, error) {
	owner, err := os.repo.GetById(id)
	if err != nil {
		os.log.WithError(err).Error("Failed to get owner with id ", id)
		return nil, errors.New("failed to fetch owner")
	}

	ownerResponse := mapper.ToOwnerResponse(owner)
	return ownerResponse, nil
}

func (os *OwnerServiceImpl) BulkDeleteOwners(ids []string, force bool) error {
	for _, id := range ids {
		owner, err := os.repo.GetById(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("owner not found: " + id)
			}
			return err
		}

		if !force {
			hasAssociations := os.repo.OwnerHasAssociations(id)
			if hasAssociations {
				return errors.New("owner " + owner.CMSUserName + " has associated records")
			}
		}

		if force {
			if err := os.repo.ForceDeleteOwner(id); err != nil {
				return err
			}
		} else {
			if err := os.repo.DeleteOwnerByID(id); err != nil {
				return err
			}
		}
	}
	return nil
}
