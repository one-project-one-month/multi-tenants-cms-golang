package authentication

import (
	"context"
	"errors"
	authenticationpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/authentication"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"sync"
)

func (s *Service) Register(
	ctx context.Context,
	req *authenticationpb.RegisterRequest,
) (*authenticationpb.RegisterResponse, error) {
	org, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "missing organization context")
	}

	orgValues := org.Get("x-organisation")
	var orgName string
	if len(orgValues) > 0 {
		orgName = orgValues[0]
	} else {
		return nil, status.Error(codes.InvalidArgument, "organization header required")
	}

	s.logger.WithFields(logrus.Fields{
		"method": "Register",
		"org":    orgName,
	}).Info("organization found")

	if err := s.validateRegisterRequest(req); err != nil {
		return nil, err
	}

	if err := s.checkRateLimit(ctx, req.GetEmail()); err != nil {
		return nil, status.Error(codes.ResourceExhausted, "rate limit exceeded for registration. Please wait 30 mins")
	}

	s.logger.WithFields(logrus.Fields{
		"method": "Register",
		"email":  req.GetEmail(),
	}).Info("processing registration request")

	flow, err := s.store.WholeRegistrationFlow(s.databaseCtx, req, orgName)
	if err != nil {
		if errors.Is(err, errors.New("namespace does not exist")) {
			return nil, status.Error(codes.NotFound, "namespace not found")
		}
	}
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("failed to process registration request")
		return nil, err
	}

	code, _ := s.generateToken(6)

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
