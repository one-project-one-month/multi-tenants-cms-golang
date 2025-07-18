package authentication

import (
	"context"
	"encoding/json"
	"github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	"time"

	"github.com/multi-tenants-cms-golang/lms-sys/internal/types"
	cook "github.com/multi-tenants-cms-golang/lms-sys/pkg/utils/cookies"
	authenticationpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/authentication"
	"github.com/sirupsen/logrus"
)

type AuthenticationInterface interface {
	buildLoginCookies(user *types.UserInfo, tokens *types.TokenPair) []string
	generateSessionID() string
	generateCSRFToken() string
	buildLogoutCookies() []string
	Register(ctx context.Context, req *authenticationpb.RegisterRequest) (*authenticationpb.RegisterResponse, error)
	sendEmailVerificationCode(email string, code string) error
	VerifyEmail(ctx context.Context, req *authenticationpb.EmailVerifyRequest) (*authenticationpb.EmailVerifyResponse, error)
	processEmailVerification(ctx context.Context, email string) error
	generateAndStoreToken(email string) error
	generateToken(length int) (string, error)
	sendVerificationEmail(email string, token string) error
	validateRegisterRequest(req *authenticationpb.RegisterRequest) error
	validateEmailVerifyRequest(req *authenticationpb.EmailVerifyRequest) error
	isValidEmail(email string) bool
	isValidPassword(password string) bool
	checkRateLimit(ctx context.Context, email string) error
	checkResendRateLimit(ctx context.Context, email string) error
	ResendVerificationEmail(ctx context.Context, req *authenticationpb.ResendVerificationRequest) (*authenticationpb.ResendVerificationResponse, error)
}

type AuthenticationService struct {
	authenticationpb.UnimplementedAuthenticationServiceServer
	databaseCtx context.Context
	store       repo.Store
	logger      *logrus.Logger
	cfg         *types.Config
	//redisClient *redis.Client
	//natsConn    *nats.Conn
}

func NewAuthenticationService(
	store repo.Store,
	logger *logrus.Logger,
	config *types.Config,
) *AuthenticationService {
	return &AuthenticationService{
		databaseCtx: context.Background(),
		store:       store,
		logger:      logger,
		cfg:         config,
		//redisClient: redisClient,
		//natsConn:    natsConn,
	}
}
func (s *AuthenticationService) buildLoginCookies(user *types.UserInfo, tokens *types.TokenPair) []string {
	var cookies []string

	cookies = append(cookies, cook.NewCookieBuilder("access_token", tokens.AccessToken).
		MaxAge(900).
		Path("/").
		Build())

	cookies = append(cookies, cook.NewCookieBuilder("refresh_token", tokens.RefreshToken).
		MaxAge(604800).
		Path("/auth").
		Build())

	userInfo := types.UserInfo{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}
	userInfoJSON, _ := json.Marshal(userInfo)
	cookies = append(cookies, cook.NewCookieBuilder("user_info", string(userInfoJSON)).
		HttpOnly(false).
		MaxAge(3600).
		Build())

	csrfToken := s.generateCSRFToken()
	cookies = append(cookies, cook.NewCookieBuilder("csrf_token", csrfToken).
		HttpOnly(false).
		MaxAge(3600).
		Build())

	sessionID := s.generateSessionID()
	cookies = append(cookies, cook.NewCookieBuilder("session_id", sessionID).
		MaxAge(7200).
		Build())

	preferences := `{"theme":"light","language":"en"}`
	cookies = append(cookies, cook.NewCookieBuilder("user_preferences", preferences).
		HttpOnly(false).
		MaxAge(2592000).
		Build())

	return cookies
}

func (s *AuthenticationService) generateSessionID() string {
	return generateRandomString(32)
}

func (s *AuthenticationService) generateCSRFToken() string {
	return generateRandomString(32)
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		time.Sleep(time.Nanosecond) // Ensure different seed for each iteration
	}
	return string(b)
}

func (s *AuthenticationService) buildLogoutCookies() []string {
	expiredTime := time.Now().Add(-24 * time.Hour)

	cookieNames := []string{
		"access_token",
		"refresh_token",
		"user_info",
		"csrf_token",
		"session_id",
		"user_preferences",
	}

	var cookies []string
	for _, name := range cookieNames {
		cookie := cook.NewCookieBuilder(name, "").
			Expires(expiredTime).
			MaxAge(0).
			Build()
		cookies = append(cookies, cookie)
	}

	return cookies
}
