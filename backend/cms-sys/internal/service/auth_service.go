package service

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/repository"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/types"
	"github.com/multi-tenants-cms-golang/cms-sys/pkg/utils"
	"github.com/nats-io/nats.go"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"image/png"
	"math/big"
	"sync"
	"time"
)

type AuthService interface {
	Login(email, password string) (*types.AuthResponse, error)
	Register(req *types.RegisterRequest) (*types.AuthResponse, error)
	RefreshToken(refreshToken string) (*types.TokenResponse, error)
	GetUserProfile(userID uuid.UUID) (*types.UserResponse, error)
	GetMe(req types.GetMeRequest) (*types.UserResponse, error)
	UpdateUserProfile(id uuid.UUID, req types.UserUpdateRequest) (*types.UserResponse, error)
	Logout(accessToken, refreshToken string) error
	GenerateMFATokenSecret(userID uuid.UUID) (*types.MFASetupResponse, error)
	VerifyMFALogin(userID uuid.UUID, verificationCode string) (*types.AuthResponse, error)
	VerifyMFASetup(userID uuid.UUID, tokenID uint, verificationCode string) error
	VerifyCredentials(email string, password string) (*types.CMSUser, error)
	EmailServiceCommunication(userId uuid.UUID) error
	VerifyEmail(email string, code string) error
}

type Service struct {
	log         *logrus.Logger
	repo        repository.AuthRepository
	redisClient *redis.Client
	natsConn    *nats.Conn
}

var _ AuthService = (*Service)(nil)

func NewService(
	log *logrus.Logger,
	repo repository.AuthRepository,
	redisClient *redis.Client,
	natsConn *nats.Conn,
) *Service {
	return &Service{
		log:         log,
		repo:        repo,
		redisClient: redisClient,
		natsConn:    natsConn,
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

	userResponse := types.UserResponse{
		ID:        user.CMSUserID,
		Name:      user.CMSUserName,
		Email:     user.CMSUserEmail,
		Role:      user.CMSUserRole,
		Verified:  user.Verified,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	err = s.EmailServiceCommunication(user.CMSUserID)
	if err != nil {
		s.log.WithError(err).Error("Failed to send email")
		return nil, err
	}
	return &types.AuthResponse{
		User: userResponse,
		////AccessToken:  accessToken,
		////RefreshToken: refreshToken,
		//ExpiresAt:    time.Now().Add(15 * time.Minute),
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

func (s *Service) GetMe(req types.GetMeRequest) (*types.UserResponse, error) {

	user, err := s.repo.GetUserByEmail(req.Email)

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

func (s *Service) UpdateUserProfile(id uuid.UUID, req types.UserUpdateRequest) (*types.UserResponse, error) {

	user, err := s.repo.GetUserByID(id)
	if err != nil {
		s.log.WithError(err).Error("Failed to get user profile with id ", id)
		return nil, errors.New("user not found")
	}

	user.CMSUserName = req.Name
	if err := s.repo.UpdateUser(user); err != nil {
		s.log.WithError(err).Error("Failed to update user profile")
		return nil, errors.New("failed to update user profile")
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

func (s *Service) Logout(accessToken, refreshToken string) error {
	if accessToken != "" {
		accessClaims, err := utils.ValidateToken(accessToken)
		if err == nil && accessClaims.ID != "" {
			ttl := time.Until(accessClaims.ExpiresAt.Time)
			if ttl > 0 {
				if err := utils.RevokeToken(accessClaims.ID, ttl); err != nil {
					s.log.WithError(err).Error("Failed to revoke access token")
					return errors.New("failed to revoke access token")
				}
			}
		}
	}

	if refreshToken != "" {
		refreshClaims, err := utils.ValidateToken(refreshToken)
		if err == nil && refreshClaims.ID != "" {
			ttl := time.Until(refreshClaims.ExpiresAt.Time)
			if ttl > 0 {
				if err := utils.RevokeToken(refreshClaims.ID, ttl); err != nil {
					s.log.WithError(err).Error("Failed to revoke refresh token")
					return errors.New("failed to revoke refresh token")
				}
			}
		}
	}

	return nil
}

func (s *Service) GenerateMFATokenSecret(userID uuid.UUID) (*types.MFASetupResponse, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		s.log.WithError(err).Error("Failed to get user profile")
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "cms-sys",
		AccountName: user.CMSUserEmail,
		Algorithm:   otp.AlgorithmSHA1,
		Digits:      otp.DigitsSix,
		Period:      30,
	})
	if err != nil {
		s.log.WithError(err).Error("Failed to generate TOTP secret")
		return nil, errors.New("failed to generate TOTP secret")
	}

	mfaToken := &types.MFAToken{
		UserID:    userID,
		MFAToken:  key.Secret(),
		ExpiresAt: func() *time.Time { t := time.Now().Add(10 * time.Minute); return &t }(),
		CreatedAt: time.Now(),
	}
	if err := s.repo.CreateMFAToken(mfaToken); err != nil {
		s.log.WithError(err).Error("Failed to generate TOTP secret")
		return nil, errors.New("failed to generate TOTP secret")
	}
	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		s.log.WithError(err).Error("Failed to generate TOTP secret")
		return nil, errors.New("failed to generate TOTP secret")
	}
	err = png.Encode(&buf, img)
	if err != nil {
		s.log.WithError(err).Error("Failed to generate TOTP secret")
		return nil, errors.New("failed to generate TOTP secret")
	}

	return &types.MFASetupResponse{
		Secret:      key.Secret(),
		QRCodeURL:   key.URL(),
		QRCodeImage: base64.StdEncoding.EncodeToString(buf.Bytes()),
		TokenID:     mfaToken.TokenID,
		ManualEntry: key.Secret(),
	}, nil

}

func (s *Service) VerifyMFASetup(userID uuid.UUID, tokenID uint, verificationCode string) error {
	mfaToken, err := s.repo.GetMFAToken(tokenID, userID)
	if err != nil {
		return errors.New("failed to get mfa token")
	}

	if mfaToken.ExpiresAt != nil && time.Now().After(*mfaToken.ExpiresAt) {
		s.log.Debug("token expired")
		return errors.New("mfa token expired")
	}
	valid := totp.Validate(verificationCode, mfaToken.MFAToken)
	if !valid {

		return errors.New("mfa token is invalid")
	}
	err = s.repo.UpdateUserMFAStatus(mfaToken.UserID, true)
	if err != nil {
		s.log.WithError(err).Error("Failed to update mfa status")
		return errors.New("failed to update mfa status")
	}

	mfaToken.ExpiresAt = nil
	if err := s.repo.UpdateMFAToken(mfaToken); err != nil {
		s.log.WithError(err).Error("Failed to update mfa token")
		return errors.New("failed to update mfa token")
	}
	return nil
}
func (s *Service) VerifyMFALogin(userID uuid.UUID, verificationCode string) (*types.AuthResponse, error) {
	mfaToken, err := s.repo.GetActiveMFAToken(userID)
	if err != nil {
		s.log.WithError(err).Error("Failed to get active mfa token")
		return nil, errors.New("failed to get mfa token")
	}
	valid := totp.Validate(verificationCode, mfaToken.MFAToken)
	if !valid {
		return nil, errors.New("mfa token is invalid")
	}
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		s.log.WithError(err).Error("Failed to get user profile")
		return nil, errors.New("failed to get user profile")
	}
	accessToken, err := utils.GenerateAccessToken(user.CMSUserID, user.CMSUserEmail, user.CMSUserRole)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}
	refreshToken, err := utils.GenerateRefreshToken(user.CMSUserID, user.CMSUserEmail, user.CMSUserRole)
	if err != nil {
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
func (s *Service) EmailServiceCommunication(userId uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	code, err := s.generateEmailVerificationCode()
	if err != nil {
		s.log.WithError(err).Error("Failed to generate verification code")
		return err
	}

	user, err := s.repo.GetUserByID(userId)
	if err != nil {
		s.log.WithError(err).Error("Failed to get user details")
		return err
	}

	var redisErr, natsErr error
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		redisErr = s.redisClient.Set(ctx, userId.String(), code, time.Minute*10).Err()
		if redisErr != nil {
			s.log.WithError(redisErr).Error("Failed to save verification code to redis")
		}
	}()

	go func() {
		defer wg.Done()
		emailPayload := EmailMessage{
			UserID:    user.CMSUserID,
			Email:     user.CMSUserEmail,
			Name:      user.CMSUserName,
			Code:      code,
			Type:      "email_verification",
			ExpiresAt: time.Now().Add(time.Minute * 10),
			Timestamp: time.Now(),
		}

		jsonPayload, err := json.Marshal(emailPayload)
		if err != nil {
			s.log.WithError(err).Error("Failed to marshal email payload")
			natsErr = err
			return
		}

		natsErr = utils.PublishMessage("email.verification", jsonPayload)
		if natsErr != nil {
			s.log.WithError(natsErr).Error("Failed to send email verification message")
		}
	}()

	wg.Wait()

	if redisErr != nil {
		s.log.WithError(redisErr).Error("Redis operation failed")
		return fmt.Errorf("failed to store verification code: %w", redisErr)
	}

	if natsErr != nil {
		s.log.WithError(natsErr).Error("NATS operation failed")
		return fmt.Errorf("failed to send email verification: %w", natsErr)
	}

	s.log.Info("Email verification code sent successfully")
	return nil
}
func (s *Service) generateEmailVerificationCode() (string, error) {
	maxNumber := big.NewInt(9999999)
	n, err := rand.Int(rand.Reader, maxNumber)
	if err != nil {
		return "", err
	}
	code := n.Int64() + 100000
	numberCode := fmt.Sprintf("%06d", code)
	return numberCode, nil
}

func (s *Service) IsMFAEnabled(userID uuid.UUID) (bool, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return false, err
	}
	return user.MFAEnabled, nil
}

func (s *Service) VerifyCredentials(email string, password string) (*types.CMSUser, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := utils.CheckPassword(password, user.Password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *Service) VerifyEmail(email string, code string) error {
	ctx := context.Background()
	userInDb, err := s.repo.GetUserByEmail(email)

	if err != nil {
		return errors.New("user not found")
	}
	storedCode, err := s.redisClient.Get(ctx, userInDb.CMSUserID.String()).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return errors.New("verification code expired")
		}
		return errors.New("failed to verify code")
	}

	if storedCode != code {
		return errors.New("invalid verification code")
	}

	if err := s.repo.UpdateUserVerificationStatus(userInDb.CMSUserID, true); err != nil {
		return errors.New("failed to update verification status")
	}

	s.redisClient.Del(ctx, userInDb.CMSUserID.String())

	return nil
}

type EmailMessage struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	Type      string    `json:"type"`
	ExpiresAt time.Time `json:"expires_at"`
	Timestamp time.Time `json:"timestamp"`
}
