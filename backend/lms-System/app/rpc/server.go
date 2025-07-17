package rpc

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"runtime/debug"
	"syscall"

	"github.com/golang-jwt/jwt/v5"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/metadata"
	"github.com/multi-tenants-cms-golang/lms-sys/app/gateway/middleware"
	authSrv "github.com/multi-tenants-cms-golang/lms-sys/app/rpc/authentication"
	db "github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	"github.com/multi-tenants-cms-golang/lms-sys/internal/types"
	authpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/authentication"
	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

const (
	component = "lms-grpc-server"
	grpcAddr  = ":9001"
	httpAddr  = ":9002"
)

type Server struct {
	store       *db.Store
	logger      *logrus.Logger
	jwtSecret   string
	jwtIssuer   string
	jwtAudience string
}

func NewServer(
	db *db.Store,
	logger *logrus.Logger,
	jwtSecret string,
	jwtIssuer string,
	jwtAudience string,
) *Server {
	return &Server{
		store:       db,
		logger:      logger,
		jwtSecret:   jwtSecret,
		jwtIssuer:   jwtIssuer,
		jwtAudience: jwtAudience,
	}
}

func (s *Server) Run() error {
	// Setup tracing
	exporter, err := stdout.New(stdout.WithPrettyPrint())
	if err != nil {
		s.logger.Errorf("failed to initialize exporter: %v", err)
		return err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	defer func() { _ = tp.Shutdown(context.Background()) }()

	// Setup metrics
	srvMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
		),
		grpcprom.WithContextLabels("tenant_name"),
	)
	reg := prometheus.NewRegistry()
	reg.MustRegister(srvMetrics)

	panicsTotal := promauto.With(reg).NewCounter(prometheus.CounterOpts{
		Name: "grpc_req_panics_recovered_total",
		Help: "Total number of gRPC requests recovered from internal panic.",
	})

	// Create interceptor chain
	interceptors := []grpc.UnaryServerInterceptor{
		// Metrics interceptor
		srvMetrics.UnaryServerInterceptor(
			grpcprom.WithExemplarFromContext(s.exemplarFromContext),
			grpcprom.WithLabelsFromContext(s.labelsFromContext),
		),
		// Logging interceptor
		logging.UnaryServerInterceptor(s.interceptorLogger(), logging.WithFieldsFromContext(s.logTraceID)),

		selector.UnaryServerInterceptor(
			auth.UnaryServerInterceptor(s.authFunction()),
			selector.MatchFunc(s.authMatcher),
		),
		// Recovery interceptor
		recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(s.panicRecoveryHandler(panicsTotal))),
	}

	// Create gRPC server
	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(interceptors...),
	)

	// Register services
	authService := authSrv.NewAuthenticationService(s.store, s.logger, &types.Config{})
	authpb.RegisterAuthenticationServiceServer(grpcServer, authService)
	srvMetrics.InitializeMetrics(grpcServer)

	// Setup run group for graceful shutdown
	var g run.Group

	// gRPC server
	g.Add(func() error {
		listener, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			return fmt.Errorf("failed to listen: %v", err)
		}
		s.logger.Infof("Starting gRPC server on %s", grpcAddr)
		return grpcServer.Serve(listener)
	}, func(err error) {
		grpcServer.GracefulStop()
	})

	// HTTP metrics server
	g.Add(func() error {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		}))
		httpSrv := &http.Server{Addr: httpAddr, Handler: mux}
		s.logger.Infof("Starting metrics server on %s", httpAddr)
		return httpSrv.ListenAndServe()
	}, func(error) {
		// HTTP server shutdown handled by run group
	})

	// Signal handler
	g.Add(run.SignalHandler(context.Background(), syscall.SIGINT, syscall.SIGTERM))

	return g.Run()
}

// Helper functions

func (s *Server) interceptorLogger() logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		s.logger.WithFields(logrus.Fields{}).Log(logrus.Level(lvl), msg)
	})
}

func (s *Server) logTraceID(ctx context.Context) logging.Fields {
	if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
		return logging.Fields{"traceID", span.TraceID().String()}
	}
	return nil
}

func (s *Server) exemplarFromContext(ctx context.Context) prometheus.Labels {
	if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
		return prometheus.Labels{"traceID": span.TraceID().String()}
	}
	return nil
}

func (s *Server) labelsFromContext(ctx context.Context) prometheus.Labels {
	labels := prometheus.Labels{}
	md := metadata.ExtractIncoming(ctx)
	if tenantName := md.Get("tenant-name"); tenantName != "" {
		labels["tenant_name"] = tenantName
	} else {
		labels["tenant_name"] = "unknown"
	}
	return labels
}

func (s *Server) authFunction() func(ctx context.Context) (context.Context, error) {
	return func(ctx context.Context) (context.Context, error) {
		token, err := auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			return nil, err
		}

		userCtx, err := s.verifyToken(token)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		return context.WithValue(ctx, middleware.UserContextKey, userCtx), nil
	}
}

func (s *Server) authMatcher(ctx context.Context, callMeta interceptors.CallMeta) bool {
	if healthpb.Health_ServiceDesc.ServiceName == callMeta.Service {
		return false
	}

	serviceName, methodName := parseFullMethod(callMeta.Method)
	if serviceName == "lms.authentication.AuthenticationService" {
		switch methodName {
		case "Register", "Login", "VerifyEmail", "ResendVerificationEmail":
			return false
		}
	}
	return true
}

func (s *Server) panicRecoveryHandler(panicsTotal prometheus.Counter) func(p any) (err error) {
	return func(p any) (err error) {
		panicsTotal.Inc()
		s.logger.Errorf("recovered from panic: %v\n%s", p, debug.Stack())
		return status.Errorf(codes.Internal, "%s", p)
	}
}

func (s *Server) verifyToken(tokenString string) (*middleware.UserContext, error) {
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

	if err := s.validateClaims(claims); err != nil {
		return nil, err
	}

	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("missing or invalid subject claim")
	}

	rolesClaim, ok := claims["roles"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("missing or invalid roles claim")
	}

	var roles []middleware.Role
	for _, r := range rolesClaim {
		if roleStr, ok := r.(string); ok {
			roles = append(roles, middleware.Role(roleStr))
		}
	}

	if len(roles) == 0 {
		return nil, fmt.Errorf("no valid roles found in token")
	}

	return &middleware.UserContext{
		UserID: userID,
		Roles:  roles,
	}, nil
}

func (s *Server) validateClaims(claims jwt.MapClaims) error {
	iss, err := claims.GetIssuer()
	if err != nil || iss != s.jwtIssuer {
		return fmt.Errorf("invalid issuer")
	}

	aud, err := claims.GetAudience()
	if err != nil || !contains(aud, s.jwtAudience) {
		return fmt.Errorf("invalid audience")
	}

	if _, err := claims.GetExpirationTime(); err != nil {
		return fmt.Errorf("invalid expiration")
	}

	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
