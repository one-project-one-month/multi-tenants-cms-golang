package authentication

import (
	"context"
	"fmt"
	gc "github.com/multi-tenants-cms-golang/lms-sys/internal/convert/global"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/cookies"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/utils"
	authenticationpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/authentication"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"time"
)

func (s *Service) Login(
	ctx context.Context,
	req *authenticationpb.LoginRequest,
) (
	*authenticationpb.LoginResponse,
	error,
) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md["x-organisation"]) == 0 {
		return nil, status.Error(codes.InvalidArgument, "X-Organisations not provided")
	}

	// Get organization from metadata
	organization := md["x-organisation"][0]

	userPasswordInDB, err := s.store.GetUserPasswordMfaByNameSpaceDomain(ctx, *gc.ConvertStringToPgText(req.GetNamespaceDomainEmail()))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	err = utils.CheckPassword(req.Password, userPasswordInDB.Password)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "permission denied. Invalid username or password")
	}

	accessToken, err := utils.GenerateAccessToken(
		userPasswordInDB.LmsUserID,
		req.GetNamespaceDomainEmail(),
		organization,
		[]byte("sK6faZaOTk2HldN9pcou/r2jCkDcT4UvaC1chMbzCuI="),
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate access token")
	}

	refreshToken, err := utils.GenerateRefreshToken(
		userPasswordInDB.LmsUserID,
		req.GetNamespaceDomainEmail(),
		organization,
		[]byte("sK6faZaOTk2HldN9pcou/r2jCkDcT4UvaC1chMbzCuI="),
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate refresh token")
	}

	accessTokenCookie := cookies.CreateSecureCookie("access_token", accessToken, time.Now().Add(time.Minute*30), "/", true)
	refreshTokenCookie := cookies.CreateSecureCookie("refresh_token", refreshToken, time.Now().Add(time.Hour*164), "/", true)

	userInfo := fmt.Sprintf(`{"user_id":"%s","email":"%s","organization":"%s"}`,
		userPasswordInDB.LmsUserID, req.GetNamespaceDomainEmail(), organization)
	userInfoCookie := cookies.CreateSecureCookie("user_info", userInfo, time.Now().Add(time.Minute*30), "/", false)
	csrfToken := s.generateCSRFToken()
	csrfCookie := cookies.CreateSecureCookie("csrf_token", csrfToken, time.Now().Add(time.Hour*164), "/", false)

	header := metadata.Pairs(
		"access-token-cookie", accessTokenCookie,
		"refresh-token-cookie", refreshTokenCookie,
		"user-info-cookie", userInfoCookie,
		"csrf-token-cookie", csrfCookie,
	)

	err = grpc.SendHeader(ctx, header)
	if err != nil {
		s.logger.Error("failed to send grpc header" + err.Error())
		return nil, err
	}

	return &authenticationpb.LoginResponse{
		Message:   "Login successful",
		MfaEnable: userPasswordInDB.MfaEnable.Bool,
	}, nil
}
