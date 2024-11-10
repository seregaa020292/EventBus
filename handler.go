package eventbus

import "context"

type (
	HandlerID = uint16
	Handlers  = map[HandlerID]Handler
	Handler   interface {
		Handle(ctx context.Context, payload any)
	}
)

type HandlerFunc func(ctx context.Context, payload any)

func (h HandlerFunc) Handle(ctx context.Context, payload any) {
	h(ctx, payload)
}
