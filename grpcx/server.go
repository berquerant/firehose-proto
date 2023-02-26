package grpcx

import (
	"context"
	"fmt"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Logger(ctx context.Context) *zap.Logger {
	return ctxzap.Extract(ctx)
}

func NewServer(srv *grpc.Server, port int) *Server {
	reflection.Register(srv)
	return &Server{
		port: port,
		srv:  srv,
	}
}

type Server struct {
	port int
	srv  *grpc.Server
}

func (s *Server) Start() error {
	ls, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen %d %w", s.port, err)
	}
	if err := s.srv.Serve(ls); err != nil {
		return fmt.Errorf("serve error %w", err)
	}
	return nil
}

func (s *Server) Stop() {
	s.srv.GracefulStop()
}
