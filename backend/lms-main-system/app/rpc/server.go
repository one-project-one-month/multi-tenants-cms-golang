package rpc

import (
	"context"
	"fmt"
	"github.com/hibiken/asynq"
	middleware2 "github.com/multi-tenants-cms-golang/lms-sys/app/middleware"
	db "github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/infra/ipfs"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/mailer"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/utils"
	"google.golang.org/grpc/reflection"
	"net"
	"syscall"
	"time"

	"github.com/golang-jwt/jwt/v5"
	authSrv "github.com/multi-tenants-cms-golang/lms-sys/app/rpc/authentication"
	lsv "github.com/multi-tenants-cms-golang/lms-sys/app/rpc/lesson"
	msv "github.com/multi-tenants-cms-golang/lms-sys/app/rpc/module"
	qsv "github.com/multi-tenants-cms-golang/lms-sys/app/rpc/quiz"
	"github.com/multi-tenants-cms-golang/lms-sys/internal/types"
	authpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/authentication"
	lspb "github.com/multi-tenants-cms-golang/lms-sys/protogen/lesson"
	mspb "github.com/multi-tenants-cms-golang/lms-sys/protogen/modules"
	qspb "github.com/multi-tenants-cms-golang/lms-sys/protogen/quiz"
	"github.com/oklog/run"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

//const (
//	grpcAddr = ":9001"
//)

type Server struct {
	store          db.Store
	logger         *logrus.Logger
	jwtSecret      string
	jwtIssuer      string
	jwtAudience    string
	grpcAddr       string
	mfaManager     *utils.MFAManager
	authMiddleware *middleware2.AuthMiddleware
	ipfsClient     *ipfs.Client
	asynqClient    *asynq.Client
	mailer         *mailer.Mailer
}

func NewServer(
	db db.Store,
	logger *logrus.Logger,
	jwtSecret string,
	jwtIssuer string,
	jwtAudience string,
	grpcAddr string,
	mfaManager *utils.MFAManager,
	authMiddleware *middleware2.AuthMiddleware,
	ipfsClient *ipfs.Client,
	asynqClient *asynq.Client,
	mailer *mailer.Mailer,
) *Server {
	return &Server{
		store:          db,
		logger:         logger,
		jwtSecret:      jwtSecret,
		jwtIssuer:      jwtIssuer,
		jwtAudience:    jwtAudience,
		grpcAddr:       grpcAddr,
		mfaManager:     mfaManager,
		authMiddleware: authMiddleware,
		ipfsClient:     ipfsClient,
		asynqClient:    asynqClient,
		mailer:         mailer,
	}
}

func (s *Server) Run() error {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware2.MetadataLoggerInterceptor(s.logger)),
		//grpc.UnaryInterceptor(s.authMiddleware.UnaryServerInterceptor()),
	)

	authService := authSrv.NewAuthenticationService(s.store, s.logger, &types.Config{
		BaseURL:           "http://localhost:9001",
		TokenLength:       0,
		TokenTTL:          0,
		MaxRetries:        1000,
		RetryDelay:        0,
		RateLimitAttempts: 100,
		RateLimitWindow:   time.Second * 20,
		Email: struct {
			NoReplyAddress string
			Templates      struct {
				Verification string
				Footer       string
			}
		}{
			NoReplyAddress: "swanhtetaunpg@gmail.com",
			Templates: struct {
				Verification string
				Footer       string
			}{
				Verification: "templates/verification",
				Footer:       "templates/footer",
			},
		},
	},
		s.mfaManager,
		s.mailer,
	)

	//fileService := files.NewFileService(s.ipfsClient, s.asynqClient, s.logger, &task.TaskCreator{})
	moduleService := msv.NewModuleService(s.store, s.logger)
	authpb.RegisterAuthenticationServiceServer(grpcServer, authService)
	mspb.RegisterModuleServiceServer(grpcServer, moduleService)

	// Lesson
	lessonService := lsv.NewLessonService(s.store, s.logger)
	lspb.RegisterLessonServiceServer(grpcServer, lessonService)

	// Quiz
	quizService := qsv.NewQuizService(s.store, s.logger)
	qspb.RegisterQuizServiceServer(grpcServer, quizService)

	//fb.RegisterFileServiceServer(grpcServer, fileService)
	var g run.Group

	g.Add(func() error {
		reflection.Register(grpcServer)
		listener, err := net.Listen("tcp", s.grpcAddr)
		if err != nil {
			return fmt.Errorf("failed to listen: %v", err)
		}
		s.logger.Infof("Starting gRPC server on %s", s.grpcAddr)
		return grpcServer.Serve(listener)
	}, func(err error) {
		grpcServer.GracefulStop()
	})

	g.Add(run.SignalHandler(context.Background(), syscall.SIGINT, syscall.SIGTERM))

	return g.Run()
}

func (s *Server) verifyToken(tokenString string) (*middleware2.UserContext, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("token parsing failed: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	//
	//// Validate standard claims
	//if err := claims.VerifyIssuer(s.jwtIssuer, true); err != nil {
	//	return nil, fmt.Errorf("invalid issuer: %w", err)
	//}
	//
	//if err := claims.VerifyAudience(s.jwtAudience, true); err != nil {
	//	return nil, fmt.Errorf("invalid audience: %w", err)
	//}

	// Extract user info
	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("missing or invalid subject claim")
	}

	rolesClaim, ok := claims["roles"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("missing or invalid roles claim")
	}

	var roles []middleware2.Role
	for _, r := range rolesClaim {
		if roleStr, ok := r.(string); ok {
			roles = append(roles, middleware2.Role(roleStr))
		}
	}

	if len(roles) == 0 {
		return nil, fmt.Errorf("no valid roles found in token")
	}

	return &middleware2.UserContext{
		UserID: userID,
		Roles:  roles,
	}, nil
}
