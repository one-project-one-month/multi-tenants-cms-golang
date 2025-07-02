package service

import (
	"errors"
	"github.com/google/uuid"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/repository"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/types"
	"github.com/multi-tenants-cms-golang/cms-sys/pkg/utils"
	"github.com/sirupsen/logrus"
	"time"
)

type AuthService interface {
	Login(email, password string) (*types.AuthResponse, error)
	Register(req *types.RegisterRequest) (*types.AuthResponse, error)
	RefreshToken(refreshToken string) (*types.TokenResponse, error)
	GetUserProfile(userID uuid.UUID) (*types.UserResponse, error)
}

type Service struct {
	log  *logrus.Logger
	repo repository.AuthRepository
}

var _ AuthService = (*Service)(nil)

func NewService(log *logrus.Logger, repo repository.AuthRepository) *Service {
	return &Service{
		log:  log,
		repo: repo,
	}
}

func (s *Service) Login(email, password string) (*types.AuthResponse, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		s.log.WithError(err).Error("Failed to get user by email")
		return nil, errors.New("invalid credentials")
	}

	if err := utils.CheckPassword(password, user.Password); err != nil {
		s.log.WithError(err).Error("Invalid password")
		return nil, errors.New("invalid credentials")
	}

	accessToken, err := utils.GenerateAccessToken(user.CMSUserID, user.CMSUserEmail, user.CMSUserRole)
	if err != nil {
		s.log.WithError(err).Error("Failed to generate access token")
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := utils.GenerateRefreshToken(user.CMSUserID, user.CMSUserEmail, user.CMSUserRole)
	if err != nil {
		s.log.WithError(err).Error("Failed to generate refresh token")
		return nil, errors.New("failed to generate refresh token")
	}

	userResponse := types.UserResponse{
		ID:        user.CMSUserID,
		Name:      user.CMSUserName,
		Email:     user.CMSUserEmail,
		Role:      user.CMSUserRole,
		Verified:  user.Verified,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return &types.AuthResponse{
		User:         userResponse,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(15 * time.Minute),
	}, nil
}

func (s *Service) Register(req *types.RegisterRequest) (*types.AuthResponse, error) {
	exists, err := s.repo.EmailExists(req.Email)
	if err != nil {
		s.log.WithError(err).Error("Failed to check if email exists")
		return nil, errors.New("failed to check email availability")
	}

	if exists {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		s.log.WithError(err).Error("Failed to hash password")
		return nil, errors.New("failed to process password")
	}

	role := req.Role
	if role == "" {
		role = string(types.CMSCustomer)
	}

	user := &types.CMSUser{
		CMSUserID:    uuid.New(),
		CMSUserName:  req.Name,
		CMSUserEmail: req.Email,
		Password:     hashedPassword,
		CMSUserRole:  role,
		Verified:     false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.CreateUser(user); err != nil {
		s.log.WithError(err).Error("Failed to create user")
		return nil, errors.New("failed to create user")
	}

	accessToken, err := utils.GenerateAccessToken(user.CMSUserID, user.CMSUserEmail, user.CMSUserRole)
	if err != nil {
		s.log.WithError(err).Error("Failed to generate access token")
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := utils.GenerateRefreshToken(user.CMSUserID, user.CMSUserEmail, user.CMSUserRole)
	if err != nil {
		s.log.WithError(err).Error("Failed to generate refresh token")
		return nil, errors.New("failed to generate refresh token")
	}

	userResponse := types.UserResponse{
		ID:        user.CMSUserID,
		Name:      user.CMSUserName,
		Email:     user.CMSUserEmail,
		Role:      user.CMSUserRole,
		Verified:  user.Verified,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return &types.AuthResponse{
		User:         userResponse,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(15 * time.Minute),
	}, nil
}

func (s *Service) RefreshToken(refreshToken string) (*types.TokenResponse, error) {
	claims, err := utils.ValidateToken(refreshToken)
	if err != nil {
		s.log.WithError(err).Error("Invalid refresh token")
		return nil, errors.New("invalid refresh token")
	}

	if claims.TokenType != "refresh" {
		return nil, errors.New("invalid token type")
	}

	user, err := s.repo.GetUserByID(claims.UserID)
	if err != nil {
		s.log.WithError(err).Error("Failed to get user")
		return nil, errors.New("user not found")
	}

	newAccessToken, err := utils.GenerateAccessToken(user.CMSUserID, user.CMSUserEmail, user.CMSUserRole)
	if err != nil {
		s.log.WithError(err).Error("Failed to generate new access token")
		return nil, errors.New("failed to generate access token")
	}

	newRefreshToken, err := utils.GenerateRefreshToken(user.CMSUserID, user.CMSUserEmail, user.CMSUserRole)
	if err != nil {
		s.log.WithError(err).Error("Failed to generate new refresh token")
		return nil, errors.New("failed to generate refresh token")
	}

	return &types.TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    time.Now().Add(15 * time.Minute),
	}, nil
}

func (s *Service) GetUserProfile(userID uuid.UUID) (*types.UserResponse, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		s.log.WithError(err).Error("Failed to get user profile")
		return nil, errors.New("user not found")
	}

	return &types.UserResponse{
		ID:        user.CMSUserID,
		Name:      user.CMSUserName,
		Email:     user.CMSUserEmail,
		Role:      user.CMSUserRole,
		Verified:  user.Verified,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
