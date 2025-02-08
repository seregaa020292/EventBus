package eventbus

import "sync"

type Middleware func(next Handler) Handler

type middleware struct {
	middlewares []Middleware
}

func (m *chainMiddlewares) Append(mw ...Middleware) {
	m.middlewares = append(m.middlewares, mw...)
}

func (m *chainMiddlewares) Wrap(h Handler) Handler {
	for i := len(m.middlewares) - 1; i >= 0; i-- {
		h = m.middlewares[i](h)
	}
	return h
}
