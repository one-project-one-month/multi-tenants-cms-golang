package authentication

import (
	"context"
	"errors"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/utils"
	authenticationpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/authentication"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
	"sync"
)

func (s *Service) Register(
	ctx context.Context,
	req *authenticationpb.RegisterRequest,
) (*authenticationpb.RegisterResponse, error) {
	org, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, utils.ErrMissingOrganization().ToGRPCStatus()
	}

	orgValues := org.Get("x-organisation")
	var orgName string
	if len(orgValues) > 0 {
		orgName = orgValues[0]
	} else {
		return nil, utils.ErrMissingOrganization().ToGRPCStatus()
	}

	s.logger.WithFields(logrus.Fields{
		"method": "Register",
		"org":    orgName,
	}).Info("organization found")

	if err := s.validateRegisterRequest(req); err != nil {
		return nil, err
	}

	if err := s.checkRateLimit(ctx, req.GetEmail()); err != nil {
		return nil, utils.ErrRateLimitExceeded("registration", "30 minutes").ToGRPCStatus()
	}

	s.logger.WithFields(logrus.Fields{
		"method": "Register",
		"email":  req.GetEmail(),
	}).Info("processing registration request")
	flow, err := s.store.WholeRegistrationFlow(s.databaseCtx, req, orgName)
	if err != nil {
		if errors.Is(err, errors.New("namespace does not exist")) {
			return nil, utils.ErrOrganizationNotFound(orgName).ToGRPCStatus()
		}
		if errors.Is(err, errors.New("user already exists")) {
			return nil, utils.ErrUserAlreadyExists(req.GetEmail()).ToGRPCStatus()
		}

		s.logger.WithError(err).Error("registration flow failed")
		return nil, utils.ErrDatabaseOperation("register user", err).ToGRPCStatus()
	}

	code, err := s.generateToken(6)
	if err != nil {
		s.logger.WithError(err).Error("failed to generate verification code")
		return nil, utils.ErrInternal("failed to generate verification code").ToGRPCStatus()
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := s.sendEmailVerificationCode(req.GetEmail(), code, orgName); err != nil {
			s.logger.WithFields(logrus.Fields{
				"error": err.Error(),
				"email": req.GetEmail(),
			}).Error("failed to send email verification code")
		}
	}()

	wg.Wait()
	return flow, nil
}
