package server

import "google.golang.org/grpc/test/bufconn"

// GrpcServer contract
type GrpcServer interface {
	Run() error
	RunMock() (*bufconn.Listener, error)
}
