package envx

import (
	"context"
	"os"
	"os/signal"

	"github.com/caarlos0/env/v7"
	"go.uber.org/zap"
)

//go:generate go run github.com/berquerant/goconfig@v0.2.0 -field "Context context.Context|Signals []os.Signal"  -option -output runner_config_generated.go

type Runner[T any] struct {
	config T
}

// NewRunner returns a new [Runner].
// Requires the appropreate struct for [Parse].
//
// [Parse]: https://pkg.go.dev/github.com/caarlos0/env/v7#Parse
func NewRunner[T any](config T) *Runner[T] {
	return &Runner[T]{
		config: config,
	}
}

// Run invokes f with config loaded from environment variables and context that is marked done when signal alives.
//
// Options:
//   - WithContext: pass the parent context, default is [context.Background].
//   - WithSignals: signals to be detected, default is [os.Interrupt].
func (r *Runner[T]) Run(
	f func(context.Context, T) error,
	opt ...ConfigOption,
) error {
	config := NewConfigBuilder().
		Context(context.Background()).
		Signals([]os.Signal{os.Interrupt}).
		Build()
	config.Apply(opt...)

	ctx, stop := signal.NotifyContext(
		config.Context.Get(),
		config.Signals.Get()...,
	)
	defer stop()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	opts := env.Options{
		OnSet: func(tag string, value any, isDefault bool) {
			logger.Info("Read env",
				zap.String("name", tag),
				zap.Any("value", value),
				zap.Bool("default", isDefault),
			)
		},
	}

	if err := env.Parse(&r.config, opts); err != nil {
		logger.Fatal("Failed to read env",
			zap.Error(err),
		)
		return err
	}

	return f(ctx, r.config)
}
