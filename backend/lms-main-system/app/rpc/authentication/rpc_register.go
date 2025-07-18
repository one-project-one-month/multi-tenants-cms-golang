package authentication

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/multi-tenants-cms-golang/lms-sys/pkg/utils/nats"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/utils/redis"
	authenticationpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/authentication"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (s *AuthenticationService) Register(
	ctx context.Context,
	req *authenticationpb.RegisterRequest,
) (*authenticationpb.RegisterResponse, error) {
	org, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "missing organization context")
	}

	orgValues := org.Get("x-organisation")
	var orgName string
	if len(orgValues) > 0 {
		orgName = orgValues[0]
	} else {
		return nil, status.Error(codes.InvalidArgument, "organization header required")
	}

	s.logger.WithFields(logrus.Fields{
		"method": "Register",
		"org":    orgName,
	}).Info("organization found")

	if err := s.validateRegisterRequest(req); err != nil {
		return nil, err
	}

	if err := s.checkRateLimit(ctx, req.GetEmail()); err != nil {
		return nil, status.Error(codes.ResourceExhausted, "rate limit exceeded for registration. Please wait 30 mins")
	}

	s.logger.WithFields(logrus.Fields{
		"method": "Register",
		"email":  req.GetEmail(),
	}).Info("processing registration request")

	flow, err := s.store.WholeRegistrationFlow(s.databaseCtx, req, orgName)

	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("failed to process registration request")
		return nil, err
	}

	code, _ := s.generateToken(6)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := s.sendEmailVerificationCode(req.GetEmail(), code, orgName); err != nil {
			s.logger.WithFields(logrus.Fields{
				"error": err.Error(),
				"email": req.GetEmail(),
			}).Error("failed to send email verification code")
		}
	}()

	wg.Wait()
	return flow, nil
}

func (s *AuthenticationService) sendEmailVerificationCode(email, code, organization string) error {
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

func (s *AuthenticationService) sendPasswordResetCode(email, code, organization string) error {
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

func (s *AuthenticationService) sendWelcomeEmail(email, name, organization string) error {
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

func (s *AuthenticationService) VerifyEmail(
	ctx context.Context,
	req *authenticationpb.EmailVerifyRequest,
) (*authenticationpb.EmailVerifyResponse, error) {
	if err := s.validateEmailVerifyRequest(req); err != nil {
		return nil, err
	}

	storedToken, err := redis.GetRedis("verify-token:" + req.GetEmail())
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"email": req.GetEmail(),
		}).Error("failed to retrieve token from Redis")
		return nil, status.Error(codes.NotFound, "verification token not found or expired")
	}

	if storedToken == "" {
		return nil, status.Error(codes.NotFound, "verification token not found or expired")
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

	if err := redis.DeleteRedis("verify-token:" + req.GetEmail()); err != nil {
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
