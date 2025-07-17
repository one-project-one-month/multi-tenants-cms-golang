package authentication

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/utils"
	convertor "github.com/multi-tenants-cms-golang/lms-sys/pkg/utils/convert"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/utils/redis"
	authenticationpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/authentication"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
	"time"
)

const (
	DefaultTokenLength      = 32
	TokenTTL                = 15 * time.Minute
	MaxRegistrationAttempts = 3
	RateLimitWindow         = time.Hour
)

func toRegisterSql(reqModel *authenticationpb.RegisterRequest) *repo.RegisterUserWithRolesParams {
	return &repo.RegisterUserWithRolesParams{
		LmsUserName:  reqModel.Username,
		LmsUserEmail: reqModel.Email,
		Password:     reqModel.Password,
		Address: pgtype.Text(sql.NullString{
			String: reqModel.Address,
			Valid:  reqModel.Address == "",
		}),
		PhoneNumber: pgtype.Text(sql.NullString{
			String: reqModel.PhoneNumber,
			Valid:  true,
		}),
	}
}
func (s *AuthenticationService) Register(
	ctx context.Context,
	req *authenticationpb.RegisterRequest,
) (*authenticationpb.RegisterResponse, error) {
	org, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		s.logger.WithFields(logrus.Fields{
			"method": "Register",
		}).Info("missing context")
	}
	get := org.Get("x-organization")
	s.logger.WithFields(logrus.Fields{
		"method": "Register",
		"org":    get,
	}).Info("organization found")
	if err := s.validateRegisterRequest(req); err != nil {
		return nil, err
	}

	if err := s.checkRateLimit(ctx, req.GetEmail()); err != nil {
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"method": "Register",
		"email":  req.GetEmail(),
	}).Info("processing registration request")
	modelToBeCreated := toRegisterSql(req)
	userCreated, err := s.store.RegisterUserWithRoles(s.databaseCtx, *modelToBeCreated)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"email": req.GetEmail(),
		}).Error("failed to register user")

		if strings.Contains(err.Error(), "duplicate key") {
			return nil, status.Error(codes.AlreadyExists, "user with this email already exists")
		}

		return nil, status.Error(codes.Internal, "failed to create user account")
	}

	go func() {
		if err := s.processEmailVerification(context.Background(), req.Email); err != nil {
			s.logger.WithFields(logrus.Fields{
				"error": err.Error(),
				"email": req.Email,
			}).Error("failed to process email verification")
		}
	}()

	return &authenticationpb.RegisterResponse{
		RegisteredUser: &authenticationpb.SystemUser{
			Id:               userCreated.LmsUserID.String(),
			Username:         userCreated.LmsUserName,
			Email:            userCreated.LmsUserEmail,
			PhoneNumber:      convertor.NullableStringToString(sql.NullString(userCreated.PhoneNumber)),
			Address:          convertor.NullableStringToString(sql.NullString(userCreated.Address)),
			RegistrationDate: utils.ParseTimestamp(userCreated.RegistrationDate),
			CreatedAt:        utils.ParseTimestamp(userCreated.CreatedAt),
			UpdatedAt:        utils.ParseTimestamp(userCreated.UpdatedAt),
		},
		Message: "user created successfully. email is sent to the registered mail",
	}, nil
}

func (s *AuthenticationService) VerifyEmail(
	ctx context.Context,
	req *authenticationpb.EmailVerifyRequest,
) (*authenticationpb.EmailVerifyResponse, error) {

	if err := s.validateEmailVerifyRequest(req); err != nil {
		return nil, err
	}

	storedToken, err := redis.GetRedis(req.GetEmail())
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"email": req.GetEmail(),
		}).Error("failed to retrieve token from Redis")
		return nil, status.Error(codes.NotFound, "verification token not found or expired")
	}

	if storedToken != req.GetToken() {
		s.logger.WithFields(logrus.Fields{
			"email": req.GetEmail(),
		}).Warn("invalid verification token provided")
		return nil, status.Error(codes.InvalidArgument, "invalid verification token")
	}

	err = s.store.UpdateEmailVerification(s.databaseCtx, req.GetEmail())
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"email": req.GetEmail(),
		}).Error("failed to update email verification status")
		return nil, status.Error(codes.Internal, "failed to verify email")
	}

	if err := redis.DeleteRedis(req.GetEmail()); err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"email": req.GetEmail(),
		}).Warn("failed to clean up verification token")
	}

	s.logger.WithFields(logrus.Fields{
		"email": req.GetEmail(),
	}).Info("email verification successful")

	return &authenticationpb.EmailVerifyResponse{
		Message: "Email verified successfully",
	}, nil
}
