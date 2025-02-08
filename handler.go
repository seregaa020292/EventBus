package eventbus

import (
	"context"
	"sync"
)

type HandlerID = uint64

type Handler interface {
	Handle(ctx context.Context, event Event) error
}

type HandlerFunc func(ctx context.Context, event Event) error

func (f HandlerFunc) Handle(ctx context.Context, event Event) error {
	return f(ctx, event)
}

type ErrHandler func(error)

type handler struct {
	async bool
	base  Handler
}

type HandlerOption func(*handler)

func WithHandlerAsync() HandlerOption {
	return func(h *handler) {
		h.async = true
	}
}

type Middleware func(next Handler) Handler

type chainMiddlewares struct {
	middlewares []Middleware
	mu          sync.RWMutex
}

func (m *chainMiddlewares) Append(mw ...Middleware) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.middlewares = append(m.middlewares, mw...)
}

func (m *chainMiddlewares) Wrap(h Handler) Handler {
	m.mu.RLock()
	defer m.mu.RUnlock()

	next := h
	for i := len(m.middlewares) - 1; i >= 0; i-- {
		next = m.middlewares[i](next)
	}
	return next
}
