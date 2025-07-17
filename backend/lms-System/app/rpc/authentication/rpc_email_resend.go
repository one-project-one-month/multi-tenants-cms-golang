package authentication

import (
	"context"
	authenticationpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/authentication"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AuthenticationService) ResendVerificationEmail(
	ctx context.Context,
	req *authenticationpb.ResendVerificationRequest,
) (*authenticationpb.ResendVerificationResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if !s.isValidEmail(req.GetEmail()) {
		return nil, status.Error(codes.InvalidArgument, "invalid email format")
	}
	user, err := s.store.GetUserByEmail(s.databaseCtx, req.GetEmail())
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	if user.EmailVerified.Valid {
		return nil, status.Error(codes.FailedPrecondition, "email already verified")
	}

	if err := s.checkResendRateLimit(ctx, req.GetEmail()); err != nil {
		return nil, err
	}

	if err := s.processEmailVerification(ctx, req.GetEmail()); err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"email": req.GetEmail(),
		}).Error("failed to resend verification email")
		return nil, status.Error(codes.Internal, "failed to resend verification email")
	}

	return &authenticationpb.ResendVerificationResponse{
		Message: "Verification email sent successfully",
		Success: true,
	}, nil
}
