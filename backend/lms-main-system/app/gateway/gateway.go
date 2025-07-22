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
	"github.com/multi-tenants-cms-golang/lms-sys/app/gateway/ipfs"
	"github.com/multi-tenants-cms-golang/lms-sys/app/gateway/modifier"
	stripe "github.com/multi-tenants-cms-golang/lms-sys/app/gateway/stripe"
	_ "github.com/multi-tenants-cms-golang/lms-sys/doc/statik"
	authenticationpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/authentication"
	fb "github.com/multi-tenants-cms-golang/lms-sys/protogen/files"
	mspb "github.com/multi-tenants-cms-golang/lms-sys/protogen/modules"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	strp "github.com/stripe/stripe-go/v76"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

type Gateway struct {
	logger              *logrus.Logger
	grpcAddr            string
	httpAddr            string
	redisAddr           string
	swaggerDir          string
	stripeSecret        string
	stripeWebhookSecret string
}

func NewGateway(
	logger *logrus.Logger,
	grpcAddr, httpAddr, redisAddr string,
) *Gateway {
	return &Gateway{
		logger:              logger,
		grpcAddr:            grpcAddr,
		httpAddr:            httpAddr,
		redisAddr:           redisAddr,
		swaggerDir:          "../doc/swagger",
		stripeSecret:        "sk_test_51RZyq9Insa2xt1880OWZrSR1AABDcd07U38WHMO8EWBynxl0Hhv0nPJYcH8UbHWI3odWFsIZk4BchIlrinUFF6Zr00U8K2LNt8",
		stripeWebhookSecret: "whsec_94dc3e2d6e1de2cdf0d26d97a288bb3269600f42e2fe3c77c12b253474b5b17c",
	}
}

func (g *Gateway) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	jsonOption := runtime.WithMarshalerOption(
		runtime.MIMEWildcard,
		&runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				EmitUnpopulated: true,
			},
		},
	)
	strp.Key = g.stripeSecret
	gwMux := runtime.NewServeMux(
		runtime.WithErrorHandler(g.errorHandler),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{}),
		runtime.WithIncomingHeaderMatcher(g.headerMatcher),
		runtime.WithMetadata(modifier.RequestModifier),
		runtime.WithForwardResponseOption(modifier.ResponseModifier),
		jsonOption,
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
		return fmt.Errorf("failed to register authentication handler: %w", err)
	}

	err = mspb.RegisterModuleServiceHandlerFromEndpoint(
		ctx,
		gwMux,
		g.grpcAddr,
		opts,
	)
	if err != nil {
		return fmt.Errorf("failed to register module handler: %w", err)
	}

	err = fb.RegisterFileServiceHandlerFromEndpoint(
		ctx,
		gwMux,
		g.grpcAddr,
		opts,
	)
	if err != nil {
		return fmt.Errorf("failed to register application handler: %w", err)
	}
	statikFS, err := fs.New()
	if err != nil {
		return fmt.Errorf("statik filesystem error: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(statikFS)))
	mux.HandleFunc("/healthz", g.healthCheck)
	mux.Handle("/", gwMux)
	ipfsGateway := ipfs.NewIPFSGateway(g.logger, g.redisAddr)
	ipfsGateway.RegisterHandlers(gwMux)
	stripeHandler := stripe.NewStripeHandler(g.stripeSecret, g.stripeWebhookSecret, []string{
		"US",
		"GB",
		"CA",
	}, "usd")
	stripeHandler.RegisterHandlers(gwMux)
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "X-Request-ID"},
		ExposedHeaders:   []string{"X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler(mux)

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
