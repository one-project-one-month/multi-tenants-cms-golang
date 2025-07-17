package authentication

import (
	"context"

	"github.com/multi-tenants-cms-golang/lms-sys/pkg/utils/nats"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/utils/redis"
	authenticationpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/authentication"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"sync"
	"time"
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
	org, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		s.logger.WithFields(logrus.Fields{
			"method": "Register",
		}).Info("missing context")
	}
	get := org.Get("x-organization")
	s.logger.WithFields(logrus.Fields{
		"method": "Register",
		"org":    get,
	}).Info("organization found")
	if err := s.validateRegisterRequest(req); err != nil {
		return nil, err
	}

	if err := s.checkRateLimit(ctx, req.GetEmail()); err != nil {
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"method": "Register",
		"email":  req.GetEmail(),
	}).Info("processing registration request")
	nameSpace := org.Get("x-organization")[0]
	flow, err := s.store.WholeRegistrationFlow(s.databaseCtx, req, nameSpace)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Info("failed to process registration request")
		return nil, err
	}
	var wg sync.WaitGroup
	code, _ := s.generateToken(6)
	wg.Add(2)
	go func() {
		defer wg.Done()
		err := s.sendEmailVerificationCode(req.GetEmail(), code)
		if err != nil {
			s.logger.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("failed to send email verification code")
			return
		}
	}()
	go func() {
		defer wg.Done()
		payload := map[string]interface{}{
			"title":        "Email Verification From LMS",
			"to":           req.GetEmail(),
			"code":         code,
			"organisation": org,
		}
		err := nats.Publish("lms.email.verification.code", payload)
		if err != nil {
			s.logger.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("failed to send email verification code")
		}
	}()
	wg.Wait()
	return flow, nil
}

func (s *AuthenticationService) sendEmailVerificationCode(email, code string) error {
	return redis.SetRedis("verify-token"+email, code, time.Minute*10)

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
