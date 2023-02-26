package main_test

import (
	"context"
	"testing"
	"time"

	"github.com/berquerant/firehose-proto/empty"
	"github.com/berquerant/firehose-proto/grpcx"
	"github.com/berquerant/firehose-test/grpctest"
	"github.com/berquerant/firehose-test/tempdir"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestRun(t *testing.T) {
	dir := tempdir.New("grpcx")
	defer dir.Close()

	runner := grpctest.NewRunner(dir, grpctest.WithHealthWait(3*time.Second))
	defer runner.Close()
	assert.Nil(t, runner.Init(context.TODO()))

	const (
		requestID = "requestFromTestRun"
	)
	var (
		client = empty.NewEmptyServiceClient(runner.Conn)
		header metadata.MD
		ctx    = metadata.NewOutgoingContext(
			context.TODO(),
			metadata.Pairs(grpcx.ClientRequestIDTag, requestID),
		)
	)
	_, err := client.Ping(ctx, new(emptypb.Empty), grpc.Header(&header))
	assert.Nil(t, err, "ping")

	clientRequestID, ok := grpcx.ExtractClientRequestIDFromMetadata(header)
	assert.True(t, ok)
	assert.Equal(t, requestID, clientRequestID)

	serverRequestID, ok := grpcx.ExtractServerRequestIDFromMetadata(header)
	assert.True(t, ok)
	assert.True(t, serverRequestID != "")
	t.Logf("ServerRequestID: %s", serverRequestID)
}
