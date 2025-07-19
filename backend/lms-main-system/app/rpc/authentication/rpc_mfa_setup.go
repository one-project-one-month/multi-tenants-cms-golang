package authentication

import (
	"context"
	authenticationpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/authentication"
)

func (s *Service) SetUpMFA(
	ctx context.Context,
	request *authenticationpb.MFASetUpRequest,
) (*authenticationpb.MFASetUpResponse, error) {
	return &authenticationpb.MFASetUpResponse{}, nil
}
