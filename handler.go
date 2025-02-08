package eventbus

import "context"

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

func generateID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}
