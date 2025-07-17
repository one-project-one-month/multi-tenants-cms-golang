package gateway

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/multi-tenants-cms-golang/lms-sys/app/gateway/middleware"
	"github.com/multi-tenants-cms-golang/lms-sys/app/gateway/modifier"
	_ "github.com/multi-tenants-cms-golang/lms-sys/doc/statik"
	authenticationpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/authentication"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	//"google.golang.org/grpc"
	//"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// JWTConfig holds JWT verification configuration
type JWTConfig struct {
	SecretKey       string
	Issuer          string
	Audience        string
	TokenExpiration time.Duration
}

type GatewayTokenVerifier struct {
	logger    *logrus.Logger
	jwtConfig JWTConfig
}

func NewGatewayTokenVerifier(logger *logrus.Logger, jwtConfig JWTConfig) *GatewayTokenVerifier {
	return &GatewayTokenVerifier{
		logger:    logger,
		jwtConfig: jwtConfig,
	}
}

// VerifyToken validates JWT token and returns user context
func (v *GatewayTokenVerifier) VerifyToken(tokenString string) (*middleware.UserContext, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(v.jwtConfig.SecretKey), nil
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

	if err := v.validateClaims(claims); err != nil {
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

func (v *GatewayTokenVerifier) validateClaims(claims jwt.MapClaims) error {
	if iss, err := claims.GetIssuer(); err != nil || iss != v.jwtConfig.Issuer {
		return fmt.Errorf("invalid issuer")
	}

	if aud, err := claims.GetAudience(); err != nil || !contains(aud, v.jwtConfig.Audience) {
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

type Gateway struct {
	logger     *logrus.Logger
	grpcAddr   string
	httpAddr   string
	swaggerDir string
	authConfig *middleware.AuthConfig
	rbacConfig *middleware.RBACConfig
	jwtConfig  JWTConfig
}

func NewGateway(
	logger *logrus.Logger,
	grpcAddr, httpAddr string,
	jwtSecret string,
	jwtIssuer string,
	jwtAudience string,
) *Gateway {
	jwtConfig := JWTConfig{
		SecretKey:       jwtSecret,
		Issuer:          jwtIssuer,
		Audience:        jwtAudience,
		TokenExpiration: 24 * time.Hour,
	}

	tokenVerifier := NewGatewayTokenVerifier(logger, jwtConfig)

	authConfig := &middleware.AuthConfig{
		Logger:        logger,
		TokenVerifier: tokenVerifier,
	}

	rbacConfig := &middleware.RBACConfig{
		Logger: logger,
		RouteRoles: map[string][]middleware.Role{
			// Admin routes
			"/api/admin":           {middleware.LMSAdmin},
			"/api/system/settings": {middleware.LMSAdmin},
			"/api/users/manage":    {middleware.LMSAdmin},

			// Instructor routes
			"/api/courses/create":     {middleware.Instructor, middleware.LMSAdmin},
			"/api/courses/update":     {middleware.Instructor, middleware.LMSAdmin},
			"/api/courses/delete":     {middleware.Instructor, middleware.LMSAdmin},
			"/api/assignments/create": {middleware.Instructor, middleware.LMSAdmin},
			"/api/assignments/grade":  {middleware.Instructor, middleware.LMSAdmin},

			// Student routes
			"/api/courses/enroll":     {middleware.Student},
			"/api/assignments/submit": {middleware.Student},
			"/api/grades/view":        {middleware.Student},

			// Shared routes
			"/api/courses/list":    {middleware.Student, middleware.Instructor, middleware.LMSAdmin},
			"/api/courses/details": {middleware.Student, middleware.Instructor, middleware.LMSAdmin},
			"/api/profile":         {middleware.Student, middleware.Instructor, middleware.LMSAdmin},
		},
	}

	return &Gateway{
		logger:     logger,
		grpcAddr:   grpcAddr,
		httpAddr:   httpAddr,
		swaggerDir: "../doc/swagger",
		authConfig: authConfig,
		rbacConfig: rbacConfig,
		jwtConfig:  jwtConfig,
	}
}

func (g *Gateway) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gwMux := runtime.NewServeMux(
		runtime.WithErrorHandler(g.errorHandler),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{}),
		runtime.WithIncomingHeaderMatcher(g.headerMatcher),
		runtime.WithIncomingHeaderMatcher(runtime.DefaultHeaderMatcher),
		runtime.WithOutgoingHeaderMatcher(runtime.DefaultHeaderMatcher),
		runtime.WithForwardResponseOption(modifier.ResponseModifier),
		runtime.WithMetadata(modifier.RequestModifier),
	)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(25 * 1024 * 1024)),
	}

	err := authenticationpb.RegisterAuthenticationServiceHandlerFromEndpoint(
		ctx,
		gwMux,
		g.grpcAddr,
		opts,
	)
	if err != nil {
		return fmt.Errorf("authenticationpb.RegisterAuthenticationServiceHandlerFromEndpoint error: %w", err)
	}
	statikFS, err := fs.New()
	if err != nil {
		return fmt.Errorf("statik filesystem error: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(statikFS)))
	mux.HandleFunc("/healthz", g.healthCheck)
	mux.Handle("/", gwMux)

	handlerChain := g.authConfig.AuthMiddleware(
		//g.rbacConfig.RBACMiddleware(
		//	mux,
		//),
		mux,
	)
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "X-Request-ID"},
		ExposedHeaders:   []string{"X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler(handlerChain)

	// Configure HTTP server
	server := &http.Server{
		Addr:         g.httpAddr,
		Handler:      corsHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		g.logger.Info("Shutting down HTTP gateway...")
		if err := server.Shutdown(ctx); err != nil {
			g.logger.WithError(err).Error("HTTP gateway shutdown error")
		}
	}()

	g.logger.Infof("Starting HTTP gateway on %s (gRPC backend: %s)", g.httpAddr, g.grpcAddr)
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) && err != nil {
		return fmt.Errorf("HTTP gateway start error: %w", err)
	}

	return nil
}

func (g *Gateway) errorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	g.logger.WithError(err).Error("gateway error")
	runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, w, r, err)
}

func (g *Gateway) headerMatcher(key string) (string, bool) {
	switch key {
	case "X-Request-ID", "X-Correlation-ID", "Authorization":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}

func (g *Gateway) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		return
	}
}
