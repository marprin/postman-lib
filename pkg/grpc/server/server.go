package server

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/google/uuid"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/marprin/postman-lib/pkg/panic"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

type (
	// AuthenticationType create own type of authentication type
	AuthenticationType string
	registerSvcFunc    func(s *grpc.Server)

	// GRPCConfig is the object config for initialize the grpc instance
	GRPCConfig struct {
		Host                  string
		Port                  uint
		MetricHost            string
		MetricPort            uint
		MetricGracefulTimeout time.Duration
		UseTLS                bool
		ServerCertFile        string
		ServerKeyFile         string
		ClientKey             string
		SecretKey             string
		AuthenticationType    AuthenticationType
	}

	grpcServer struct {
		cfg                        *GRPCConfig
		registerSvcFunc            registerSvcFunc
		unaryInterceptorMiddleware grpc.UnaryServerInterceptor
	}
)

const (
	// AuthenticationTypeNone the default for unknown authentication
	AuthenticationTypeNone AuthenticationType = ""
	// AuthenticationTypeClientSecretKey is the authentication with client and secret key
	AuthenticationTypeClientSecretKey AuthenticationType = "client_secret_key"

	// AuthorizationHeader is the metadata key for authorization
	AuthorizationHeader string = "x-authorization"
)

var (
	ErrFailedExtractMetadata          = "Failed when extract metadata"
	ErrAuthorizationTokenIsNotPresent = "Authorization token is not present"
	ErrAuthorizationTokenIsNotValid   = "Authorization token is not valid"
)

// NewGrpcServer Initialize grpc instance
func NewGrpcServer(cfg *GRPCConfig, fn registerSvcFunc) GrpcServer {
	inst := &grpcServer{
		cfg:             cfg,
		registerSvcFunc: fn,
	}

	// Set the middleware
	if cfg.AuthenticationType == AuthenticationTypeClientSecretKey {
		inst.unaryInterceptorMiddleware = inst.clientSecretKeyUnaryInterceptorHandler
	} else {
		inst.unaryInterceptorMiddleware = inst.defaultUnaryInterceptorHandler
	}

	return inst
}

// CustomizeUnaryInterceptorMiddleware is the setter function so the user can customize the middleware by themselve
func (g *grpcServer) CustomizeUnaryInterceptorMiddleware(fnc grpc.UnaryServerInterceptor) {
	g.unaryInterceptorMiddleware = fnc
}

func (g *grpcServer) Run() error {
	promRegistry := prometheus.NewRegistry()
	grpcMetrics := grpcprometheus.NewServerMetrics()

	// Register the prometheus standard metrics
	promRegistry.MustRegister(grpcMetrics)

	httpAddr := fmt.Sprintf("%s:%d", g.cfg.Host, g.cfg.Port)
	metricAddr := fmt.Sprintf("%s:%d", g.cfg.MetricHost, g.cfg.MetricPort)

	httpListener, err := net.Listen("tcp", httpAddr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
			grpcMetrics.UnaryServerInterceptor(),
			g.unaryInterceptorMiddleware,
		)),
	)

	// Register the proto service
	if g.registerSvcFunc != nil {
		g.registerSvcFunc(grpcServer)
	}

	// Register the health check service
	healthpb.RegisterHealthServer(grpcServer, health.NewServer())

	// Register the reflection of grpc server
	reflection.Register(grpcServer)

	// Initialize the metrics
	grpcMetrics.InitializeMetrics(grpcServer)

	httpMetricServer := &http.Server{
		Handler: promhttp.HandlerFor(promRegistry, promhttp.HandlerOpts{}),
		Addr:    metricAddr,
	}

	logrus.Infoln("Starting Metrics Server")
	go func() {
		if err := httpMetricServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Errorf("Failed to start HTTP metric server for GRPC on %s", metricAddr)
		}
	}()

	errChan := make(chan error, 1)
	logrus.Infoln("Starting GRPC Server")
	go func() {
		if err := grpcServer.Serve(httpListener); err != nil {
			logrus.Errorf("Failed to start GRPC Server on %s", httpAddr)
			errChan <- err
		}
	}()
	logrus.Infof("GRPC server started on %s and Metrics server started on %s", httpAddr, metricAddr)

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, os.Interrupt, os.Kill)

	select {
	case sig := <-termChan:
		logrus.Infof("Receive terminating signal, prepare for shutdown service, %s", sig)

		logrus.Infoln("Trying to terminate GRPC Server")
		grpcServer.GracefulStop()
		logrus.Infoln("Successfully graceful stop GRPC server")

		logrus.Infoln("Trying to terminate metrics server")
		ctx, cancel := context.WithTimeout(context.Background(), g.cfg.MetricGracefulTimeout*time.Second)
		defer cancel()

		if err := httpMetricServer.Shutdown(ctx); err != nil {
			logrus.Infof("Force shutdown metrics server, error: %s", err)
		}
		logrus.Infoln("Successfully terminate Metric server")
	case err := <-errChan:
		pErr := fmt.Errorf("Exiting with error: %+v", err)
		return errors.New(pErr.Error())
	}

	return nil
}

func (g *grpcServer) RunMock() (*bufconn.Listener, error) {
	lis := bufconn.Listen(1024 * 1024)

	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
			grpc.UnaryServerInterceptor(g.unaryInterceptorMiddleware),
		)),
	)

	// registering services
	if g.registerSvcFunc != nil {
		g.registerSvcFunc(s)
	}
	var err error

	go func() {
		err = s.Serve(lis)
	}()

	return lis, err
}

func (g *grpcServer) clientSecretKeyUnaryInterceptorHandler(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	timeStart := time.Now()
	reqID := uuid.New().String()
	method := info.FullMethod
	// handle any panic ocured on the server
	defer panic.HandlePanic(func(r interface{}) {
		panic.ToPanicError(r, info.FullMethod)
		err = status.Errorf(codes.Internal, "%s", r)
	})

	// skip loging health check probe
	if strings.HasPrefix(method, "/grpc.health") {
		return handler(ctx, req)
	}

	// Handle the client key validation
	metadata, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, ErrFailedExtractMetadata)
	}

	if len(metadata[AuthorizationHeader]) == 0 {
		return nil, status.Error(codes.Unauthenticated, ErrAuthorizationTokenIsNotPresent)
	}

	token := metadata[AuthorizationHeader][0]
	decToken, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, ErrAuthorizationTokenIsNotValid)
	}

	clientKey := string(decToken)
	if g.cfg.ClientKey != clientKey[:len(clientKey)-1] {
		return nil, status.Error(codes.Unauthenticated, ErrAuthorizationTokenIsNotValid)
	}

	LogUnaryRequest(reqID, method, req)
	resp, err = handler(ctx, req)
	LogUnaryResponse(reqID, method, timeStart, resp, err)

	return resp, err
}

func (g *grpcServer) defaultUnaryInterceptorHandler(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	timeStart := time.Now()
	reqID := uuid.New().String()
	method := info.FullMethod
	// handle any panic ocured on the server
	defer panic.HandlePanic(func(r interface{}) {
		panic.ToPanicError(r, info.FullMethod)
		err = status.Errorf(codes.Internal, "%s", r)
	})

	// skip loging health check probe
	if strings.HasPrefix(method, "/grpc.health") {
		return handler(ctx, req)
	}

	LogUnaryRequest(reqID, method, req)
	resp, err = handler(ctx, req)
	LogUnaryResponse(reqID, method, timeStart, resp, err)

	return resp, err
}

func LogUnaryRequest(reqID, method string, req interface{}) {
	logrus.WithFields(logrus.Fields{
		"req_id": reqID,
		"method": method,
		"req":    req,
	}).Info("incoming rpc unary request")
}

func LogUnaryResponse(reqID, method string, timeStart time.Time, resp interface{}, err error) {
	fields := logrus.Fields{
		"req_id": reqID,
		"method": method,
		"took":   time.Since(timeStart),
	}

	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("rpc unary request failed")
	} else {
		fields["resp"] = resp
		logrus.WithFields(fields).Info("rpc unary request succeeded")
	}
}
