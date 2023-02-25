package empty

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	UnimplementedEmptyServiceServer
}

func (*Server) Ping(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	logger := ctxzap.Extract(ctx)
	logger.Info("Server.Ping()")
	return new(emptypb.Empty), nil
}
