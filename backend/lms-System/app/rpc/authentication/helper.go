package authentication

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/utils/nats"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/utils/redis"
	authenticationpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/authentication"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

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
		time.Sleep(100 * time.Millisecond)

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

	s.logger.WithFields(logrus.Fields{
		"email": email,
		"ttl":   ttl,
	}).Debug("verification token generated and stored")

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
		<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
			<div style="background-color: #f8f9fa; padding: 20px; border-radius: 8px;">
				<h2 style="color: #343a40; text-align: center;">Email Verification</h2>
				<p style="color: #6c757d; font-size: 16px;">
					Thank you for registering with our Learning Management System!
				</p>
				<p style="color: #6c757d; font-size: 16px;">
					Please click the button below to verify your email address:
				</p>
				<div style="text-align: center; margin: 30px 0;">
					<a href="%s" style="background-color: #007bff; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; font-weight: bold;">
						Verify Email
					</a>
				</div>
				<p style="color: #6c757d; font-size: 14px;">
					<strong>Important:</strong> This link will expire in 15 minutes.
				</p>
				<p style="color: #6c757d; font-size: 14px;">
					If you did not request this verification, please ignore this email.
				</p>
				<hr style="margin: 20px 0; border: none; border-top: 1px solid #dee2e6;">
				<p style="color: #6c757d; font-size: 12px; text-align: center;">
					This is an automated message, please do not reply.
				</p>
			</div>
		</body>
		</html>
	`, verifyURL)

	emailPayload := map[string]interface{}{
		"to":      email,
		"subject": "Email Verification - LMS System",
		"body":    emailBody,
		"type":    "verification",
	}

	return nats.Publish("lms.email.verify", emailPayload)
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

	if req.GetUsername() == "" {
		return status.Error(codes.InvalidArgument, "username is required")
	}

	if len(req.GetUsername()) < 3 {
		return status.Error(codes.InvalidArgument, "username must be at least 3 characters")
	}

	if !s.isValidPassword(req.GetPassword()) {
		return status.Error(codes.InvalidArgument, "password must contain at least one uppercase letter, one lowercase letter, and one number")
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

	if len(req.GetToken()) != DefaultTokenLength*2 {
		return status.Error(codes.InvalidArgument, "invalid token format")
	}

	return nil
}

func (s *AuthenticationService) isValidEmail(email string) bool {

	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return false
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	localPart := parts[0]
	domainPart := parts[1]

	if len(localPart) == 0 || len(domainPart) == 0 {
		return false
	}

	if len(localPart) > 64 || len(domainPart) > 253 {
		return false
	}

	return true
}

func (s *AuthenticationService) isValidPassword(password string) bool {
	hasUpper := false
	hasLower := false
	hasDigit := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		}
	}

	return hasUpper && hasLower && hasDigit
}

func (s *AuthenticationService) checkRateLimit(ctx context.Context, email string) error {
	key := fmt.Sprintf("rate_limit:register:%s", email)
	countStr, err := redis.GetRedis(key)

	maxAttempts := s.cfg.RateLimitAttempts
	if maxAttempts == 0 {
		maxAttempts = MaxRegistrationAttempts
	}

	window := s.cfg.RateLimitWindow
	if window == 0 {
		window = RateLimitWindow
	}

	var count int
	if err == nil && countStr != "" {
		count, _ = strconv.Atoi(countStr)
	}

	if count >= maxAttempts {
		s.logger.WithFields(logrus.Fields{
			"email": email,
			"count": count,
		}).Warn("rate limit exceeded for registration")
		return status.Error(codes.ResourceExhausted, "registration rate limit exceeded. Please try again later")
	}

	count++
	if err := redis.SetRedis(key, strconv.Itoa(count), window); err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"email": email,
		}).Error("failed to set rate limit counter")
		return status.Error(codes.Internal, "rate limit check failed")
	}

	s.logger.WithFields(logrus.Fields{
		"email": email,
		"count": count,
	}).Debug("rate limit check passed")

	return nil
}

func (s *AuthenticationService) checkResendRateLimit(ctx context.Context, email string) error {
	key := fmt.Sprintf("rate_limit:resend:%s", email)
	countStr, err := redis.GetRedis(key)

	maxAttempts := 3
	window := time.Hour

	var count int
	if err == nil && countStr != "" {
		count, _ = strconv.Atoi(countStr)
	}

	if count >= maxAttempts {
		s.logger.WithFields(logrus.Fields{
			"email": email,
			"count": count,
		}).Warn("rate limit exceeded for resend verification")
		return status.Error(codes.ResourceExhausted, "resend rate limit exceeded. Please try again later")
	}

	count++
	if err := redis.SetRedis(key, strconv.Itoa(count), window); err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"email": email,
		}).Error("failed to set resend rate limit counter")
		return status.Error(codes.Internal, "resend rate limit check failed")
	}

	return nil
}
