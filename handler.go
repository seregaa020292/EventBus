package eventbus

import "context"

type (
	HandlerID = uint16
	Handler   interface {
		Handle(ctx context.Context, event Event)
	}
	HandlerOption func(*handlerOption)
)

type HandlerFunc func(ctx context.Context, event Event)

func (h HandlerFunc) Handle(ctx context.Context, event Event) {
	h(ctx, event)
}

func WithHandlerIsAsync(v bool) HandlerOption {
	return func(o *handlerOption) {
		o.isAsync = v
	}
}

type (
	handlerOptions = map[HandlerID]*handlerOption
	handlerOption  struct {
		next    Handler
		isAsync bool
	}
)

func newHandlerOption(next Handler, options []HandlerOption) *handlerOption {
	handler := &handlerOption{
		next:    next,
		isAsync: false,
	}

	for _, option := range options {
		option(handler)
	}

	return handler
}

func (o handlerOption) Handle(ctx context.Context, event Event) {
	o.next.Handle(ctx, event)
}
