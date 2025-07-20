package authentication

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/multi-tenants-cms-golang/lms-sys/internal/types"
	cook "github.com/multi-tenants-cms-golang/lms-sys/pkg/cookies"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/infra/nats"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/infra/redis"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/utils"
	authenticationpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/authentication"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"strings"
	"sync"
	"time"
)

func (s *Service) processEmailVerification(ctx context.Context, email string) error {
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

func (s *Service) generateAndStoreToken(email string) error {
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

func (s *Service) generateToken(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("invalid token length: %d", length)
	}

	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return hex.EncodeToString(bytes), nil
}

func (s *Service) validateRegisterRequest(req *authenticationpb.RegisterRequest) error {
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
	result, err := utils.ParsePhoneNumberEnhanced("+48608422691", "PL")
	if err != nil {
		return status.Error(codes.Internal, "failed to parse phone number")
	}

	if !result.IsValid {
		return status.Error(codes.Internal, "phone number is invalid")
	}
	return nil
}

func (s *Service) validateEmailVerifyRequest(req *authenticationpb.EmailVerifyRequest) error {
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

func (s *Service) isValidEmail(email string) bool {

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

func (s *Service) isValidPassword(password string) bool {
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

func (s *Service) checkRateLimit(ctx context.Context, email string) error {
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

func (s *Service) checkResendRateLimit(ctx context.Context, email string) error {
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

func (s *Service) buildLoginCookies(user *types.UserInfo, tokens *types.TokenPair) []string {
	var cookies []string

	cookies = append(cookies, cook.NewCookieBuilder("access_token", tokens.AccessToken).
		MaxAge(900).
		Path("/").
		Build())

	cookies = append(cookies, cook.NewCookieBuilder("refresh_token", tokens.RefreshToken).
		MaxAge(604800).
		Path("/auth").
		Build())

	userInfo := types.UserInfo{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}
	userInfoJSON, _ := json.Marshal(userInfo)
	cookies = append(cookies, cook.NewCookieBuilder("user_info", string(userInfoJSON)).
		HttpOnly(false).
		MaxAge(3600).
		Build())

	csrfToken := s.generateCSRFToken()
	cookies = append(cookies, cook.NewCookieBuilder("csrf_token", csrfToken).
		HttpOnly(false).
		MaxAge(3600).
		Build())

	sessionID := s.generateSessionID()
	cookies = append(cookies, cook.NewCookieBuilder("session_id", sessionID).
		MaxAge(7200).
		Build())

	preferences := `{"theme":"light","language":"en"}`
	cookies = append(cookies, cook.NewCookieBuilder("user_preferences", preferences).
		HttpOnly(false).
		MaxAge(2592000).
		Build())

	return cookies
}

func (s *Service) generateSessionID() string {
	return generateRandomString(32)
}

func (s *Service) generateCSRFToken() string {
	return generateRandomString(32)
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		time.Sleep(time.Nanosecond)
	}
	return string(b)
}

func (s *Service) buildLogoutCookies() []string {
	expiredTime := time.Now().Add(-24 * time.Hour)

	cookieNames := []string{
		"access_token",
		"refresh_token",
		"user_info",
		"csrf_token",
		"session_id",
		"user_preferences",
	}

	var cookies []string
	for _, name := range cookieNames {
		cookie := cook.NewCookieBuilder(name, "").
			Expires(expiredTime).
			MaxAge(0).
			Build()
		cookies = append(cookies, cookie)
	}

	return cookies
}

func (s *Service) sendEmailVerificationCode(email, code, organization string) error {
	if err := redis.SetRedis("verify-token:"+email, code, time.Minute*10); err != nil {
		return fmt.Errorf("failed to store verification code: %w", err)
	}

	payload := map[string]any{
		"to":           email,
		"subject":      "Email Verification Code",
		"templateName": "verification_code",
		"data": map[string]interface{}{
			"code":         code,
			"organization": organization,
			"email":        email,
			"expires_at":   time.Now().Add(time.Minute * 10).Unix(),
		},
		"trackingId": fmt.Sprintf("verify_%s_%d", email, time.Now().Unix()),
	}

	if err := nats.Publish("lms.email.verification", payload); err != nil {
		return fmt.Errorf("failed to publish email verification: %w", err)
	}

	return nil
}

func (s *Service) sendPasswordResetCode(email, code, organization string) error {
	if err := redis.SetRedis("password-reset:"+email, code, time.Minute*15); err != nil {
		return fmt.Errorf("failed to store password reset code: %w", err)
	}

	payload := map[string]any{
		"type":         "password_reset",
		"email":        email,
		"code":         code,
		"organization": organization,
		"template":     "password_reset",
		"subject":      "Password Reset Code",
		"expires_at":   time.Now().Add(time.Minute * 15).Unix(),
		"sent_at":      time.Now().Unix(),
	}

	return nats.Publish("lms.email.password_reset", payload)
}

func (s *Service) sendWelcomeEmail(email, name, organization string) error {
	payload := map[string]any{
		"type":         "welcome",
		"email":        email,
		"name":         name,
		"organization": organization,
		"template":     "welcome",
		"subject":      "Welcome to LMS",
		"sent_at":      time.Now().Unix(),
	}

	return nats.Publish("lms.email.welcome", payload)
}
