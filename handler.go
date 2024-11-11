package eventbus

import "context"

type (
	HandlerID = uint16
	Handlers  = map[HandlerID]Handler
	Handler   interface {
		Handle(ctx context.Context, event Event)
	}
	HandlerOption func(*handlerOption)
)

type HandlerFunc func(ctx context.Context, event Event)

func (h HandlerFunc) Handle(ctx context.Context, event Event) {
	h(ctx, event)
}

type handlerOption struct {
	isAsync bool
	next    Handler
}

func newHandlerOption(next Handler) *handlerOption {
	return &handlerOption{
		isAsync: false,
		next:    next,
	}
}

func (o handlerOption) Handle(ctx context.Context, event Event) {
	o.next.Handle(ctx, event)
}

func WithHandlerIsAsync(s bool) HandlerOption {
	return func(o *handlerOption) {
		o.isAsync = s
	}
}
