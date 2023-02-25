package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"

	"github.com/berquerant/firehose-proto/empty"
	"github.com/berquerant/firehose-proto/grpcx"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

func main() {
	if err := func() error {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()

		s := grpc.NewServer(
			grpcmiddleware.WithUnaryServerChain(
				grpcx.NewBaseUnaryServerInterceptors()...,
			),
		)
		empty.RegisterEmptyServiceServer(s, &empty.Server{})
		r := grpcx.NewRunner(grpcx.NewServer(s, getPort()))
		r.Run(ctx)
		return r.Err()
	}(); err != nil {
		log.Fatal(err)
	}
}

func getPort() int {
	if p := os.Getenv("PORT"); p != "" {
		if i, err := strconv.Atoi(p); err == nil {
			return i
		}
	}
	return 10001
}
