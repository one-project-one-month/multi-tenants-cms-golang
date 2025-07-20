package authentication

import (
	"context"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/mailer"

	"github.com/multi-tenants-cms-golang/lms-sys/pkg/infra/redis"
	authenticationpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/authentication"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) VerifyEmail(ctx context.Context, req *authenticationpb.EmailVerifyRequest) (*authenticationpb.EmailVerifyResponse, error) {
	if err := s.validateEmailVerifyRequest(req); err != nil {
		return nil, err
	}

	if err := s.checkRateLimit(ctx, req.GetEmail()); err != nil {
		return nil, err
	}

	storedToken, err := redis.GetRedis("verify-token:" + req.GetEmail())
	if err != nil {
		s.logger.WithError(err).Error("failed to retrieve verification token")
		return nil, status.Error(codes.Internal, "verification failed")
	}

	if storedToken == "" {
		return nil, status.Error(codes.NotFound, "invalid or expired token")
	}

	if storedToken != req.GetToken() {
		return nil, status.Error(codes.InvalidArgument, "invalid token")
	}

	if err := s.store.UpdateEmailVerification(s.databaseCtx, req.GetEmail()); err != nil {
		s.logger.WithError(err).Error("failed to update verification status")
		return nil, status.Error(codes.Internal, "verification failed")
	}

	_ = redis.DeleteRedis("verify-token:" + req.GetEmail())

	return &authenticationpb.EmailVerifyResponse{
		Message: "Email verified successfully",
	}, nil
}

func (s *Service) sendVerificationEmail(email, token string) error {
	job := mailer.EmailJob{
		From:    s.cfg.Email.NoReplyAddress,
		To:      []string{email},
		Subject: "Verify Your Email Address",
		Templates: struct {
			BodyPath   string
			FooterPath string
		}{
			BodyPath:   s.cfg.Email.Templates.Verification,
			FooterPath: s.cfg.Email.Templates.Footer,
		},
		Body: map[string]interface{}{
			"Token":     token,
			"ExpiresIn": TokenTTL.Minutes(),
		},
	}

	return s.mailer.SendEmail(job)
}
