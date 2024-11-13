package eventbus

import (
	"context"
	"sync"
)

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
		wg      *sync.WaitGroup
		isAsync bool
	}
)

func newHandlerOption(next Handler, wg *sync.WaitGroup, options []HandlerOption) *handlerOption {
	handler := &handlerOption{
		next:    next,
		wg:      wg,
		isAsync: false,
	}

	for _, option := range options {
		option(handler)
	}

	return handler
}

func (h handlerOption) Handle(ctx context.Context, event Event) {
	if h.isAsync {
		h.asyncHandle(ctx, event)
		return
	}

	h.next.Handle(ctx, event)
}

func (h handlerOption) asyncHandle(ctx context.Context, event Event) {
	h.wg.Add(1)

	go func() {
		defer h.wg.Done()

		h.next.Handle(ctx, event)
	}()
}
