package grpcx

import (
	"context"
	"errors"
	"os"
	"time"

	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpcctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// WithRequestID returns a UnaryServerInterceptor that appends request-scoped metadata to context tags and grpc metadata (header).
//
// grpcx.server_request_id is a server-side request-scoped id, created when the server receives a request.
// grpcx.client_request_id is a client-side request-scoped id, may be read when the server received a request.
func WithRequestID() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		var (
			tags = grpcctxtags.Extract(ctx)
			md   = metadataFromContext(ctx)
		)
		// receive client request id
		if id, found := ExtractClientRequestIDFromMetadata(md); found {
			tags.Set(ClientRequestIDTag, id)
		}
		// create and set server request id
		serverRequestID := NewUUID()
		tags.Set(ServerRequestIDTag, serverRequestID)
		// set request ids to header
		AppendServerRequestIDIntoMetadata(md, serverRequestID)
		if err := grpc.SetHeader(ctx, md); err != nil {
			Logger(ctx).Warn("Failed to set header",
				zap.Error(err),
			)
		}

		Logger(ctx).Info("Start")
		return handler(ctx, req)
	}
}

func NewBaseUnaryServerInterceptors() []grpc.UnaryServerInterceptor {
	logger, _ := zap.NewProduction()
	grpczap.ReplaceGrpcLoggerV2(logger)
	return []grpc.UnaryServerInterceptor{
		grpcctxtags.UnaryServerInterceptor(
			grpcctxtags.WithFieldExtractor(
				grpcctxtags.CodeGenRequestFieldExtractor,
			),
		),
		grpczap.UnaryServerInterceptor(logger, grpczap.WithDurationField(func(d time.Duration) zapcore.Field {
			return zap.Int64("grpc.time_ns", d.Nanoseconds())
		})),
		WithRequestID(),
	}
}

var (
	ErrNoCredentials = errors.New("NoCredentials")
)

// NewServerCredentialsFromEnv returns a transport cred from the files specified by env.
// Returns ErrNoCredentials when FIREHOSE_PROTO_CRED_CERT or FIREHOSE_PROTO_CRED_KEY is unset.
func NewServerCredentialsFromEnv() (grpc.ServerOption, error) {
	var (
		crt = os.Getenv("FIREHOSE_PROTO_CRED_CERT")
		key = os.Getenv("FIREHOSE_PROTO_CRED_KEY")
	)
	if crt != "" && key != "" {
		cred, err := credentials.NewServerTLSFromFile(crt, key)
		if err != nil {
			return nil, err
		}
		return grpc.Creds(cred), nil
	}
	return nil, ErrNoCredentials
}

// NewClientCredentialsFromEnv returns a transport cred from the file specified by env.
// Returns an insecure cred and ErrNoCredentials when FIREHOSE_PROTO_CRED_CERT is unset.
func NewClientCredentialsFromEnv() (grpc.DialOption, error) {
	if crt := os.Getenv("FIREHOSE_PROTO_CRED_CERT"); crt != "" {
		cred, err := credentials.NewClientTLSFromFile(crt, "")
		if err != nil {
			return nil, err
		}
		return grpc.WithTransportCredentials(cred), nil
	}
	return grpc.WithTransportCredentials(insecure.NewCredentials()), ErrNoCredentials
}
