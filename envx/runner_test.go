package envx_test

import (
	"context"
	"os"
	"testing"

	"github.com/berquerant/firehose-proto/envx"
	"github.com/stretchr/testify/assert"
)

type config struct {
	World string `env:"WORLD"`
}

func TestRunner(t *testing.T) {
	const world = "NEW"
	os.Setenv("WORLD", world)

	t.Run("Run", func(t *testing.T) {
		var (
			cfg config
			run = func(_ context.Context, cfg config) error {
				assert.Equal(t, world, cfg.World)
				return nil
			}
		)
		assert.Nil(t, envx.NewRunner(cfg).Run(run))
	})

	t.Run("Cancel", func(t *testing.T) {
		var (
			cfg config
			run = func(ctx context.Context, cfg config) error {
				assert.Equal(t, world, cfg.World)
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
					return nil
				}
			}
		)

		ctx, cancel := context.WithCancel(context.TODO())
		cancel()
		assert.ErrorIs(t, context.Canceled, envx.NewRunner(cfg).Run(run, envx.WithContext(ctx)))
	})

}
