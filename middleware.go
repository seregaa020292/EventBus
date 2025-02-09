package eventbus

import "sync"

type Middleware func(next Handler) Handler

type middlewares struct {
	chain []Middleware
}

func (m *middlewares) Append(mw ...Middleware) {
	m.chain = append(m.middlewares, mw...)
}

func (m *middlewares) Wrap(h Handler) Handler {
	for i := len(m.chain) - 1; i >= 0; i-- {
		h = m.chain[i](h)
	}
	return h
}
