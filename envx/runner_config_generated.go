// Code generated by "goconfig -field Context context.Context|Signals []os.Signal -option -output runner_config_generated.go"; DO NOT EDIT.

package envx

import (
	"context"
	"os"
)

type ConfigItem[T any] struct {
	modified     bool
	value        T
	defaultValue T
}

func (s *ConfigItem[T]) Set(value T) {
	s.modified = true
	s.value = value
}
func (s *ConfigItem[T]) Get() T {
	if s.modified {
		return s.value
	}
	return s.defaultValue
}
func (s *ConfigItem[T]) Default() T {
	return s.defaultValue
}
func (s *ConfigItem[T]) IsModified() bool {
	return s.modified
}
func NewConfigItem[T any](defaultValue T) *ConfigItem[T] {
	return &ConfigItem[T]{
		defaultValue: defaultValue,
	}
}

type Config struct {
	Context *ConfigItem[context.Context]
	Signals *ConfigItem[[]os.Signal]
}
type ConfigBuilder struct {
	context context.Context
	signals []os.Signal
}

func (s *ConfigBuilder) Context(v context.Context) *ConfigBuilder {
	s.context = v
	return s
}
func (s *ConfigBuilder) Signals(v []os.Signal) *ConfigBuilder {
	s.signals = v
	return s
}
func (s *ConfigBuilder) Build() *Config {
	return &Config{
		Context: NewConfigItem(s.context),
		Signals: NewConfigItem(s.signals),
	}
}

func NewConfigBuilder() *ConfigBuilder { return &ConfigBuilder{} }
func (s *Config) Apply(opt ...ConfigOption) {
	for _, x := range opt {
		x(s)
	}
}

type ConfigOption func(*Config)

func WithContext(v context.Context) ConfigOption {
	return func(c *Config) {
		c.Context.Set(v)
	}
}
func WithSignals(v []os.Signal) ConfigOption {
	return func(c *Config) {
		c.Signals.Set(v)
	}
}
