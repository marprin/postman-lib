package client

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
)

// UnaryInterceptor is used to log the request and response of a gRPC call
func UnaryInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	elapsed := time.Since(start)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Request": req,
			"Error":   err,
			"Elapsed": elapsed,
		}).Errorln("Interceptor Log")
	} else {
		logrus.WithFields(logrus.Fields{
			"Request":  req,
			"Response": reply,
			"Elapsed":  elapsed,
		}).Infoln("Interceptor Log")
	}
	return err
}

// CreateDefaultgRPCConn is the default configuration for make the gRPC connection
func CreateDefaultgRPCConn(endpoint string, timeout time.Duration) *grpc.ClientConn {
	return CreategRPCConn(
		endpoint,
		grpc.WithInsecure(),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff:           backoff.DefaultConfig,
			MinConnectTimeout: timeout * time.Second,
		}),
		grpc.WithUnaryInterceptor(UnaryInterceptor),
	)
}

// CreategRPCConn initialize gRPC connection with user can custom the params
func CreategRPCConn(endpoint string, dialOptions ...grpc.DialOption) *grpc.ClientConn {
	cc, err := grpc.Dial(endpoint, dialOptions...)
	if err != nil {
		logrus.Error("Error when creating the gRPC connection", err)
	} else {
		logrus.Info("Successfully connect to gRPC Server: " + endpoint)
	}
	return cc
}
