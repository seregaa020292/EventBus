## Usage

```go
package main

const (
	EventTypeCreated eventbus.EventName = iota
	EventTypeUpdated
)

type ListenHandle1 struct {
	V int
}

func (h ListenHandle1) Handle(ctx context.Context, event eventbus.Event) {
	fmt.Printf("eventbus handler1 event: %+v, %d\n", event, h.V)
	//fmt.Printf("eventbus handler1 event: %+v\n", event.(eventCreated))
}

func main() {
	bus := eventbus.New()

	listenHandle1 := ListenHandle1{1}
	listenHandle2 := ListenHandle1{2}

	{
		_, unsubscribe := bus.Subscribe(EventTypeCreated, listenHandle1)
		unsubscribe()
	}
	{
		handlerID, _ := bus.Subscribe(EventTypeCreated, listenHandle2)
		bus.Unsubscribe(EventTypeCreated, handlerID)
	}

	handler2 := eventbus.HandlerFunc(func(ctx context.Context, event eventbus.Event) {
		switch e := event.(type) {
		case eventCreated:
			fmt.Printf("eventbus handler2 event: %+v, %v\n", e, e.ID)
		case eventUpdated:
			fmt.Printf("eventbus handler2 event: %+v, %v\n", e, e.ID)
		}
	})
	bus.Subscribe(EventTypeCreated, handler2)
	bus.Subscribe(EventTypeUpdated, handler2)

	bus.Publish(context.Background(), eventCreated{"Foo eventCreated"})

	//
	var events = eventbus.Events{}

	events.Enqueue(eventCreated{"Foo eventCreated12"})
	events.Enqueue(eventUpdated{"Foo eventUpdated122"})
	bus.Flush(context.Background(), &events)

	events.Enqueue(eventUpdated{"Foo eventUpdated1223"})
	bus.Flush(context.Background(), &events)

	time.Sleep(200 * time.Millisecond)
}
```
