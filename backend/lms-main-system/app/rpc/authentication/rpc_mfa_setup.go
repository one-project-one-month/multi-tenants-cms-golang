package authentication

import (
	"context"
	"fmt"
	"github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/cookies"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"time"

	authenticationpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/authentication"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) SetUpMFA(
	ctx context.Context,
	request *authenticationpb.MFASetUpRequest,
) (*authenticationpb.MFASetUpResponse, error) {
	if request.GetDomainEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	secret, err := s.mfaManager.GenerateSecret(request.GetDomainEmail())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate MFA secret")
	}

	mfaSecretID, err := s.store.SetUpMFA(ctx, repo.SetUpMFAParams{
		MfaSecret:          secret.Secret,
		LmsUserDomainEmail: request.GetDomainEmail(),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to setup MFA: %v", err))
	}
	mfaSecretCookie := cookies.CreateSecureCookie("mfa-code-id", mfaSecretID.String(), time.Now().Add(10*time.Minute), "/", true)
	header := metadata.Pairs("mfa-code", mfaSecretCookie)
	err = grpc.SendHeader(ctx, header)
	if err != nil {
		s.logger.Error("failed to send grpc header" + err.Error())
		return nil, err
	}
	return &authenticationpb.MFASetUpResponse{
		PhotoString: secret.QRCodeBase64Str,
		ManualEntry: secret.Secret,
		MfaUrl:      secret.QRCodeURL,
	}, nil
}
