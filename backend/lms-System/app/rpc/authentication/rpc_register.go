package authentication

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	convertor "github.com/multi-tenants-cms-golang/lms-sys/pkg/utils/convert"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/utils/nats"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/utils/redis"
	authenticationpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/authentication"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	DefaultTokenLength      = 32
	TokenTTL                = 15 * time.Minute
	MaxRegistrationAttempts = 3
	RateLimitWindow         = time.Hour
)

func (s *AuthenticationService) Register(
	ctx context.Context,
	req *authenticationpb.RegisterRequest,
) (*authenticationpb.RegisterResponse, error) {
	if err := s.validateRegisterRequest(req); err != nil {
		return nil, err
	}

	if err := s.checkRateLimit(ctx, req.GetEmail()); err != nil {
		return nil, err
	}

	//incomingContext, _ := metadata.FromIncomingContext(ctx)
	s.logger.WithFields(logrus.Fields{
		"method": "Register",
		"email":  req.GetEmail(),
	}).Info("processing registration request")

	modelToBeCreated := convertor.RegisterProtoToModel(req)
	user, err := s.store.RegisterLMSUser(s.databaseCtx, *modelToBeCreated)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"email": req.GetEmail(),
		}).Error("failed to register user")
		return nil, status.Error(codes.Internal, "failed to create user account")
	}

	if err := s.processEmailVerification(ctx, user.LmsUserEmail); err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"email": user.LmsUserEmail,
		}).Error("failed to process email verification")
	}

	return convertor.RegisterModelToProto(&user), nil
}

func (s *AuthenticationService) VerifyEmail(
	ctx context.Context,
	req *authenticationpb.EmailVerifyRequest,
) (*authenticationpb.EmailVerifyResponse, error) {
	if err := s.validateEmailVerifyRequest(req); err != nil {
		return nil, err
	}

	storedToken, err := redis.GetRedis(req.GetEmail())
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"email": req.GetEmail(),
		}).Error("failed to retrieve token from Redis")
		return nil, status.Error(codes.Internal, "verification failed")
	}

	if storedToken != req.GetToken() {
		s.logger.WithFields(logrus.Fields{
			"email": req.GetEmail(),
		}).Warn("invalid verification token provided")
		return nil, status.Error(codes.InvalidArgument, "invalid verification token")
	}

	err = s.store.UpdateEmailVerification(s.databaseCtx, req.GetEmail())
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"email": req.GetEmail(),
		}).Error("failed to update email verification status")
		return nil, status.Error(codes.Internal, "failed to verify email")
	}

	if err := redis.DeleteRedis(req.GetEmail()); err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"email": req.GetEmail(),
		}).Warn("failed to clean up verification token")
	}

	s.logger.WithFields(logrus.Fields{
		"email": req.GetEmail(),
	}).Info("email verification successful")

	return &authenticationpb.EmailVerifyResponse{
		Message: "Email verified successfully",
	}, nil
}

func (s *AuthenticationService) processEmailVerification(ctx context.Context, email string) error {
	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := s.generateAndStoreToken(email); err != nil {
			errChan <- fmt.Errorf("failed to generate token: %w", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		token, err := redis.GetRedis(email)
		if err != nil {
			errChan <- fmt.Errorf("failed to retrieve token for email: %w", err)
			return
		}
		if err := s.sendVerificationEmail(email, token); err != nil {
			errChan <- fmt.Errorf("failed to send verification email: %w", err)
		}
	}()

	go func() {
		wg.Wait()
		close(errChan)
	}()

	var errors []error
	for err := range errChan {
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("email verification processing failed: %v", errors)
	}

	return nil
}

func (s *AuthenticationService) generateAndStoreToken(email string) error {
	tokenLength := s.cfg.TokenLength
	if tokenLength == 0 {
		tokenLength = DefaultTokenLength
	}

	emailCode, err := s.generateToken(tokenLength)
	if err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}

	ttl := s.cfg.TokenTTL
	if ttl == 0 {
		ttl = TokenTTL
	}

	if err := redis.SetRedis(email, emailCode, ttl); err != nil {
		return fmt.Errorf("failed to store token in Redis: %w", err)
	}

	return nil
}

func (s *AuthenticationService) generateToken(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("invalid token length: %d", length)
	}

	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return hex.EncodeToString(bytes), nil
}

func (s *AuthenticationService) sendVerificationEmail(email, token string) error {
	if email == "" || token == "" {
		return fmt.Errorf("email and token are required")
	}

	baseURL := s.cfg.BaseURL
	if baseURL == "" {
		baseURL = "http://localhost:8098"
	}

	verifyURL := fmt.Sprintf("%s/verify-email?email=%s&token=%s",
		baseURL, url.QueryEscape(email), url.QueryEscape(token))

	emailBody := fmt.Sprintf(`
       <html>
       <body>
           <h2>Email Verification</h2>
           <p>Please click the link below to verify your email address:</p>
           <a href="%s">Verify Email</a>
           <p>This link will expire in 15 minutes.</p>
           <p>If you did not request this verification, please ignore this email.</p>
       </body>
       </html>
   `, verifyURL)

	return nats.Publish("lms.email.verify", []byte(emailBody))
}

func (s *AuthenticationService) validateRegisterRequest(req *authenticationpb.RegisterRequest) error {
	if req == nil {
		return status.Error(codes.InvalidArgument, "request is required")
	}

	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if !s.isValidEmail(req.GetEmail()) {
		return status.Error(codes.InvalidArgument, "invalid email format")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	if len(req.GetPassword()) < 8 {
		return status.Error(codes.InvalidArgument, "password must be at least 8 characters")
	}

	return nil
}

func (s *AuthenticationService) validateEmailVerifyRequest(req *authenticationpb.EmailVerifyRequest) error {
	if req == nil {
		return status.Error(codes.InvalidArgument, "request is required")
	}

	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if !s.isValidEmail(req.GetEmail()) {
		return status.Error(codes.InvalidArgument, "invalid email format")
	}

	if req.GetToken() == "" {
		return status.Error(codes.InvalidArgument, "token is required")
	}

	return nil
}

func (s *AuthenticationService) isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func (s *AuthenticationService) checkRateLimit(ctx context.Context, email string) error {
	key := fmt.Sprintf("rate_limit:register:%s", email)
	count, err := redis.GetRedis(key)
	if err != nil {
		return status.Error(codes.Internal, "rate limit check failed")
	}

	maxAttempts := s.cfg.RateLimitAttempts
	if maxAttempts == 0 {
		maxAttempts = MaxRegistrationAttempts
	}

	if count != "" {
		s.logger.WithFields(logrus.Fields{
			"email": email,
			"count": count,
		}).Info("checking rate limit")
	}

	window := s.cfg.RateLimitWindow
	if window == 0 {
		window = RateLimitWindow
	}

	if err := redis.SetRedis(key, "1", window); err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"email": email,
		}).Error("failed to set rate limit counter")
	}

	return nil
}
