package grpc

import (
	service "adflow/internal/ports/grpc/service"

	"google.golang.org/grpc"
)

func NewGRPCServer(a service.AdServiceServer) *grpc.Server {
	s := grpc.NewServer()
	service.RegisterAdServiceServer(s, a)
	return s
}
