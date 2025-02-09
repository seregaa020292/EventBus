package eventbus

import "sync"

type Middleware func(next Handler) Handler

type middlewares struct {
	chain []Middleware
	mu    sync.RWMutex
}

func (m *middlewares) Append(mw ...Middleware) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.chain = append(m.chain, mw...)
}

func (m *middlewares) Wrap(h Handler) Handler {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for i := len(m.chain) - 1; i >= 0; i-- {
		h = m.chain[i](h)
	}
	return h
}
