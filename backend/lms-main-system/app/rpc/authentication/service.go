package authentication

import (
	"context"
	"github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	"github.com/multi-tenants-cms-golang/lms-sys/internal/types"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/utils"
	authenticationpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/authentication"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	DefaultTokenLength      = 32
	TokenTTL                = 10 * time.Minute
	MaxRegistrationAttempts = 3
	RateLimitWindow         = 30 * time.Minute
)

type (
	AuthenticationInterface interface {
		// Cookie/session related
		buildLoginCookies(user *types.UserInfo, tokens *types.TokenPair) []string
		buildLogoutCookies() []string
		generateSessionID() string
		generateCSRFToken() string

		// Register Auth API handlers (gRPC methods)
		Register(ctx context.Context, req *authenticationpb.RegisterRequest) (*authenticationpb.RegisterResponse, error)
		VerifyEmail(ctx context.Context, req *authenticationpb.EmailVerifyRequest) (*authenticationpb.EmailVerifyResponse, error)
		ResendVerificationEmail(ctx context.Context, req *authenticationpb.ResendVerificationRequest) (*authenticationpb.ResendVerificationResponse, error)
		Login(ctx context.Context, req *authenticationpb.LoginRequest) (*authenticationpb.LoginResponse, error)
		LogOut(ctx context.Context, req *authenticationpb.LoginOutRequest) (*authenticationpb.LogoutResponse, error)
		Refresh(ctx context.Context, req *authenticationpb.RefreshTokenRequest) (*authenticationpb.RefreshTokenResponse, error)
		// Email & Token logic
		sendEmailVerificationCode(email, code, organisation string) error
		sendVerificationEmail(email string, token string) error
		generateAndStoreToken(email string) error
		generateToken(length int) (string, error)

		// Validation logic
		validateRegisterRequest(req *authenticationpb.RegisterRequest) error
		validateEmailVerifyRequest(req *authenticationpb.EmailVerifyRequest) error
		isValidEmail(email string) bool
		isValidPassword(password string) bool

		// Rate limiting
		checkRateLimit(ctx context.Context, email string) error
		checkResendRateLimit(ctx context.Context, email string) error

		// Email processing
		processEmailVerification(ctx context.Context, email string) error
	}
	Service struct {
		authenticationpb.UnimplementedAuthenticationServiceServer
		databaseCtx context.Context
		store       repo.Store
		logger      *logrus.Logger
		cfg         *types.Config
		mfaManager  *utils.MFAManager
	}
)

var _ AuthenticationInterface = (*Service)(nil)

func NewAuthenticationService(
	store repo.Store,
	logger *logrus.Logger,
	config *types.Config,
	mfaManager *utils.MFAManager,
) *Service {
	return &Service{
		databaseCtx: context.Background(),
		store:       store,
		logger:      logger,
		cfg:         config,
		mfaManager:  mfaManager,
	}
}
