package grpcx

import (
	"context"

	"go.uber.org/zap"
)

func NewRunner(server *Server) *Runner {
	return &Runner{
		Server: server,
	}
}

type Runner struct {
	*Server
	err error
}

func (r *Runner) Err() error {
	return r.err
}

func (r *Runner) Run(ctx context.Context) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	doneC := make(chan struct{})
	go func() {
		defer close(doneC)
		if err := r.Start(); err != nil {
			r.err = err
			logger.Error("Runner got an error",
				zap.Error(err),
			)
		}
	}()

	select {
	case <-doneC:
		logger.Warn("Runner closed")
	case <-ctx.Done():
		logger.Info("Runner is stopping")
		r.Stop()
		logger.Info("Runner stopped")
	}
}
