package gateway

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/multi-tenants-cms-golang/lms-sys/protogen/assignments"
	cpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/courses"
	tpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/tenants"
	mpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/modules"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Gateway struct {
	logger     *logrus.Logger
	grpcAddr   string
	httpAddr   string
	swaggerDir string
}

func NewGateway(logger *logrus.Logger, grpcAddr, httpAddr string) *Gateway {
	return &Gateway{
		logger:     logger,
		grpcAddr:   grpcAddr,
		httpAddr:   httpAddr,
		swaggerDir: "../doc/swagger",
	}
}

func (g *Gateway) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gwMux := runtime.NewServeMux(
		runtime.WithErrorHandler(g.errorHandler),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{}),
		runtime.WithIncomingHeaderMatcher(g.headerMatcher),
	)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(25 * 1024 * 1024)), // 25MB
	}

	if err := pb.RegisterAssignmentServiceHandlerFromEndpoint(ctx, gwMux, g.grpcAddr, opts); err != nil {
		return fmt.Errorf("register grpc gateway error: %w", err)
	}

	if err := tpb.RegisterTenantServiceHandlerFromEndpoint(ctx, gwMux, g.grpcAddr, opts); err != nil {
		return fmt.Errorf("failed to register tenant service gateway")
	}

	if err := cpb.RegisterCourseServiceHandlerFromEndpoint(ctx, gwMux, g.grpcAddr, opts); err != nil {
		return fmt.Errorf("failed to register course service gateway")
	}

	if err := mpb.RegisterModuleServiceHandlerFromEndpoint(ctx, gwMux, g.grpcAddr, opts); err != nil {
		return fmt.Errorf("failed to register module service gateway")
	}

	statikFS, err := fs.New()
	if err != nil {
		return fmt.Errorf("statik filesystem error: %w", err)
	}

	// Create main mux router
	mux := http.NewServeMux()
	mux.Handle("/", gwMux)
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(statikFS)))
	mux.HandleFunc("/healthz", g.healthCheck)

	// Configure CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler(mux)

	// Configure HTTP server
	server := &http.Server{
		Addr:         g.httpAddr,
		Handler:      corsHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
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
	case "X-Request-ID", "X-Correlation-ID":
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
