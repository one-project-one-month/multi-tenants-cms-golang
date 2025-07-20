package authentication

import (
	"context"
	authenticationpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/authentication"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) ResendVerificationEmail(ctx context.Context, req *authenticationpb.ResendVerificationRequest) (*authenticationpb.ResendVerificationResponse, error) {
	if err := s.checkResendRateLimit(ctx, req.GetEmail()); err != nil {
		return nil, err
	}

	if err := s.generateAndStoreToken(req.GetEmail()); err != nil {
		s.logger.WithError(err).Error("failed to resend verification email")
		return nil, status.Error(codes.Internal, "failed to resend verification email")
	}

	return &authenticationpb.ResendVerificationResponse{
		Success: true,
		Message: "Verification email resent",
	}, nil
}
