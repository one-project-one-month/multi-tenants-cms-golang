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
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"image/png"
	"log"
	"math/big"
	"sync"
	"time"
)

type AuthService interface {
	Login(email, password string) (*types.AuthResponse, error)
	Register(req *types.RegisterRequest, address string) (*types.AuthResponse, error)
	RefreshToken(refreshToken string) (*types.TokenResponse, error)
	GetUserProfile(userID uuid.UUID) (*types.UserResponse, error)
	GetMe(req types.GetMeRequest) (*types.UserResponse, error)
	UpdateUserProfile(id uuid.UUID, req types.UserUpdateRequest) (*types.UserResponse, error)
	Logout(accessToken, refreshToken string) error
	GenerateMFATokenSecret(userID uuid.UUID) (*types.MFASetupResponse, error)
	VerifyMFALogin(userID uuid.UUID, verificationCode string) (*types.AuthResponse, error)
	VerifyMFASetup(userID uuid.UUID, tokenID uint, verificationCode string) error
	VerifyCredentials(email string, password string) (*types.CMSUser, error)
	EmailServiceCommunication(userId uuid.UUID, info *utils.LocationInfo) error
	VerifyEmail(email string, code string) error
}

type Service struct {
	log         *logrus.Logger
	repo        repository.AuthRepository
	redisClient *redis.Client
	appStatus   string
	jwtSecret   []byte
}

var _ AuthService = (*Service)(nil)

func NewService(
	log *logrus.Logger,
	repo repository.AuthRepository,
	redisClient *redis.Client,
	jwtSecret []byte,

) *Service {
	return &Service{
		log:         log,
		repo:        repo,
		redisClient: redisClient,
		appStatus:   utils.GetEnv("CMS_APP_STATUS", "production"),
		jwtSecret:   jwtSecret,
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

	accessToken, err := utils.GenerateAccessToken(
		user.CMSUserID,
		user.CMSUserEmail,
		user.CMSUserRole,
		s.jwtSecret,
	)
	if err != nil {
		s.log.WithError(err).Error("Failed to generate access token")
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := utils.GenerateRefreshToken(
		user.CMSUserID,
		user.CMSUserEmail,
		user.CMSUserRole,
		s.jwtSecret,
	)
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

func (s *Service) Register(req *types.RegisterRequest, address string) (*types.AuthResponse, error) {
	locationInfo, err := utils.GetLocationFromIP(address)
	if err != nil {
		s.log.WithError(err).Error("Failed to get location info")
	}
	exists, err := s.repo.EmailExists(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email availability: %w", err)
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to process password: %w", err)
	}

	user := &types.CMSUser{
		CMSUserID:    uuid.New(),
		CMSUserName:  req.Name,
		CMSUserEmail: req.Email,
		Password:     hashedPassword,
		CMSUserRole:  s.determineUserRole(req.Role),
		Verified:     s.isAutoVerifyEnabled(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	if !s.isAutoVerifyEnabled() {
		if err := s.sendVerificationEmail(user.CMSUserID, locationInfo); err != nil {
			return nil, fmt.Errorf("failed to send verification email: %w", err)
		}
	}

	return &types.AuthResponse{
		User: types.UserResponse{
			ID:        user.CMSUserID,
			Name:      user.CMSUserName,
			Email:     user.CMSUserEmail,
			Role:      user.CMSUserRole,
			Verified:  user.Verified,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
}

func (s *Service) RefreshToken(refreshToken string) (*types.TokenResponse, error) {
	claims, err := utils.ValidateToken(
		refreshToken,
		s.jwtSecret,
	)
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

	newAccessToken, err := utils.GenerateAccessToken(
		user.CMSUserID,
		user.CMSUserEmail,
		user.CMSUserRole,
		s.jwtSecret,
	)
	if err != nil {
		s.log.WithError(err).Error("Failed to generate new access token")
		return nil, errors.New("failed to generate access token")
	}

	newRefreshToken, err := utils.GenerateRefreshToken(
		user.CMSUserID,
		user.CMSUserEmail,
		user.CMSUserRole,
		s.jwtSecret,
	)
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
		accessClaims, err := utils.ValidateToken(accessToken, s.jwtSecret)
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
		refreshClaims, err := utils.ValidateToken(refreshToken, s.jwtSecret)
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
		Issuer:      "cms-doc-sys",
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
	s.log.WithFields(logrus.Fields{
		"userID":           userID,
		"tokenID":          tokenID,
		"verificationCode": verificationCode,
	}).Debug("Starting MFA verification setup")

	mfaToken, err := s.repo.GetMFAToken(tokenID, userID)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"userID":  userID,
			"tokenID": tokenID,
			"error":   err.Error(),
		}).Error("Failed to get MFA token from repository")
		return errors.New("failed to get mfa token")
	}

	s.log.WithFields(logrus.Fields{
		"userID":  userID,
		"tokenID": tokenID,
	}).Debug("Successfully retrieved MFA token from repository")

	if mfaToken == nil {
		s.log.WithFields(logrus.Fields{
			"userID":  userID,
			"tokenID": tokenID,
		}).Error("MFA token is nil after successful repository call")
		return errors.New("mfa token not found")
	}

	s.log.WithFields(logrus.Fields{
		"userID":         userID,
		"tokenID":        tokenID,
		"mfaTokenUserID": mfaToken.UserID,
		"hasExpiresAt":   mfaToken.ExpiresAt != nil,
		"mfaTokenLength": len(mfaToken.MFAToken),
	}).Debug("MFA token details")

	if mfaToken.ExpiresAt != nil && time.Now().After(*mfaToken.ExpiresAt) {
		s.log.WithFields(logrus.Fields{
			"userID":      userID,
			"tokenID":     tokenID,
			"expiresAt":   mfaToken.ExpiresAt,
			"currentTime": time.Now(),
		}).Debug("MFA token has expired")
		return errors.New("mfa token expired")
	}

	s.log.WithFields(logrus.Fields{
		"userID":  userID,
		"tokenID": tokenID,
	}).Debug("MFA token expiration check passed")

	// Add nil check for MFAToken field
	if mfaToken.MFAToken == "" {
		s.log.WithFields(logrus.Fields{
			"userID":  userID,
			"tokenID": tokenID,
		}).Error("MFA token secret is empty")
		return errors.New("invalid mfa token")
	}

	s.log.WithFields(logrus.Fields{
		"userID":         userID,
		"tokenID":        tokenID,
		"mfaTokenSecret": mfaToken.MFAToken[:10] + "...",
	}).Debug("MFA token secret validation passed, validating TOTP code")

	valid := totp.Validate(verificationCode, mfaToken.MFAToken)
	if !valid {
		s.log.WithFields(logrus.Fields{
			"userID":           userID,
			"tokenID":          tokenID,
			"verificationCode": verificationCode,
		}).Error("TOTP validation failed - invalid verification code")
		return errors.New("mfa token is invalid")
	}

	s.log.WithFields(logrus.Fields{
		"userID":  userID,
		"tokenID": tokenID,
	}).Debug("TOTP validation successful, updating user MFA status")

	err = s.repo.UpdateUserMFAStatus(mfaToken.UserID, true)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"userID":         userID,
			"tokenID":        tokenID,
			"mfaTokenUserID": mfaToken.UserID,
			"error":          err.Error(),
		}).Error("Failed to update user MFA status")
		return errors.New("failed to update mfa status")
	}

	s.log.WithFields(logrus.Fields{
		"userID":  userID,
		"tokenID": tokenID,
	}).Debug("Successfully updated user MFA status, updating MFA token")

	mfaToken.ExpiresAt = nil
	if err := s.repo.UpdateMFAToken(mfaToken); err != nil {
		s.log.WithFields(logrus.Fields{
			"userID":  userID,
			"tokenID": tokenID,
			"error":   err.Error(),
		}).Error("Failed to update MFA token")
		return errors.New("failed to update mfa token")
	}

	s.log.WithFields(logrus.Fields{
		"userID":  userID,
		"tokenID": tokenID,
	}).Info("MFA setup verification completed successfully")

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
	accessToken, err := utils.GenerateAccessToken(
		user.CMSUserID,
		user.CMSUserEmail,
		user.CMSUserRole,
		s.jwtSecret,
	)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}
	refreshToken, err := utils.GenerateRefreshToken(
		user.CMSUserID,
		user.CMSUserEmail,
		user.CMSUserRole,
		s.jwtSecret,
	)
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

func (s *Service) EmailServiceCommunication(userId uuid.UUID, info *utils.LocationInfo) error {
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

	userLoc, err := time.LoadLocation(info.Timezone)
	if err != nil {
		s.log.WithError(err).Error("Failed to load timezone")
	}

	expirationDuration := time.Minute * 10
	expirationTimeLocal := time.Now().In(userLoc).Add(expirationDuration)
	expirationUnix := expirationTimeLocal.UTC()
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
		emailPayload := types.EmailRequest{
			To:           user.CMSUserEmail,
			Subject:      "verification code",
			TemplateName: "verification",
			Data: map[string]interface{}{
				"CompanyNameLogo": "CMS System",
				"CompanyName":     "CMS System",
				"Subject":         "Email verification",
				"Code":            code,
				"ExpirationTime":  expirationUnix,
				"CurrentYear":     time.Now().Year(),
			},
		}

		data, _ := json.Marshal(emailPayload)
		if err := utils.PublishMessage("email.verification", data); err != nil {
			log.Printf("Failed to publish message: %v", err)
		}

		//natsErr = utils.PublishMessage("email.verification", data)
		//if natsErr != nil {
		//	s.log.WithError(natsErr).Error("Failed to send email verification message")
		//}
	}()

	wg.Wait()

	if redisErr != nil {
		s.log.WithError(redisErr).Error("Redis operation failed")
		return fmt.Errorf("failed to store.go verification code: %w", redisErr)
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

func (s *Service) determineUserRole(requestedRole string) string {
	if requestedRole == "" {
		return string(types.CMSCustomer)
	}
	return requestedRole
}

func (s *Service) isAutoVerifyEnabled() bool {
	switch s.appStatus {
	case "development":
		return true
	case "staging", "production":
		return false
	default:
		return false
	}
}

func (s *Service) sendVerificationEmail(userID uuid.UUID, info *utils.LocationInfo) error {
	if s.appStatus == "development" {
		s.log.Info("Skipping email sending in development mode")
		return nil
	}
	return s.EmailServiceCommunication(userID, info)
}
