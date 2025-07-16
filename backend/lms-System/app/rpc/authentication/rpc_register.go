package authentication

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	convertor "github.com/multi-tenants-cms-golang/lms-sys/pkg/utils/convert"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/utils/nats"
	authenticationpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/authentication"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
)

func (s *AuthenticationService) Register(
	ctx context.Context,
	req *authenticationpb.RegisterRequest,
) (*authenticationpb.RegisterResponse, error) {

	incomingContext, _ := metadata.FromIncomingContext(ctx)
	s.logger.WithFields(logrus.Fields{
		"method": "Register",
	}).Infof("incoming metadata: %v", incomingContext)
	modelTobCreated := convertor.RegisterProtoToModel(req)
	user, err := s.store.RegisterLMSUser(s.databaseCtx, *modelTobCreated)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("failed to register user: %v", err.Error())
		return nil, err
	}

	return convertor.RegisterModelToProto(&user), nil
}

func (s *AuthenticationService) generateToken(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
func (s *AuthenticationService) sendVerificationEmail(email string, token string) error {
	verifyURL := "http://localhost:8098/verify-email?token=" + token
	body := "Click to verify: " + verifyURL
	return nats.Publish("lms.email.verify", []byte(body))
}
