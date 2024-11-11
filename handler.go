package eventbus

import "context"

type (
	HandlerID = uint16
	Handlers  = map[HandlerID]Handler
	Handler   interface {
		Handle(ctx context.Context, event Event)
	}
)

type HandlerFunc func(ctx context.Context, event Event)

func (h HandlerFunc) Handle(ctx context.Context, event Event) {
	h(ctx, event)
}
