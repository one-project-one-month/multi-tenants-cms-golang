package authentication

import (
	"context"
	authenticationpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/authentication"
)

func (s *Service) LogOut(ctx context.Context, req *authenticationpb.LoginOutRequest) (*authenticationpb.LogoutResponse, error) {
	return nil, nil
}
