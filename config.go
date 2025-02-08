package eventbus

import (
	"fmt"
	"time"
)

type Config struct {
	AsyncTimeout time.Duration
	ErrorHandler func(error)
}

func newConfig(opts ...Option) Config {
	cfg := Config{
		AsyncTimeout: 30 * time.Second,
		ErrorHandler: defaultErrorHandler,
	}

	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

type Option func(*Config)

func WithAsyncTimeout(t time.Duration) Option {
	return func(cfg *Config) {
		cfg.AsyncTimeout = t
	}
}

func WithErrorHandler(f func(error)) Option {
	return func(cfg *Config) {
		cfg.ErrorHandler = f
	}
}

func defaultErrorHandler(err error) {
	fmt.Printf("eventbus: error handling: %v\n", err)
}
