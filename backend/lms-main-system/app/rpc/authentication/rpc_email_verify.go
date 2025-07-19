package authentication

import (
	"context"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/infra/redis"
	authenticationpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/authentication"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) VerifyEmail(
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
