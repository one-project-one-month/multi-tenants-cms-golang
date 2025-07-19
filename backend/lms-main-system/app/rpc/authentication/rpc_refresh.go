package authentication

import (
	"context"
	authenticationpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/authentication"
)

func (s *Service) Refresh(
	ctx context.Context,
	req *authenticationpb.RefreshTokenRequest,
) (*authenticationpb.RefreshTokenResponse, error) {
	return nil, nil
}
