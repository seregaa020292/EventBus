package eventbus_test

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"sanbox/eventbus"
)

type testHandler struct {
	mu      sync.Mutex
	calls   []any
	handler func(ctx context.Context, payload any)
}

func (h *testHandler) Handle(ctx context.Context, payload any) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.calls = append(h.calls, payload)
	if h.handler != nil {
		h.handler(ctx, payload)
	}
}

func TestEventBus(t *testing.T) {
	bus := eventbus.New()
	ctx := context.Background()

	// Создаем обработчик
	var calls []any
	handler := &testHandler{
		handler: func(ctx context.Context, payload any) {
			calls = append(calls, payload)
		},
	}

	// Подписываемся на событие
	eventType := eventbus.EventType(1)
	_, unsubscribe := bus.Subscribe(eventType, handler)

	// Проверяем, что подписка работает
	bus.Publish(ctx, eventbus.NewEvent(eventType, "test payload"))

	// Убедимся, что обработчик был вызван
	assert.Len(t, calls, 1)
	assert.Equal(t, "test payload", calls[0])

	// Отписываемся и проверяем, что обработчик больше не вызывается
	unsubscribe()
	bus.Publish(ctx, eventbus.NewEvent(eventType, "another payload"))

	// Убедимся, что обработчик не был вызван
	assert.Len(t, calls, 1)
}

func TestFlush(t *testing.T) {
	bus := eventbus.New()
	ctx := context.Background()
	events := &eventbus.Events{}

	// Создаем обработчик
	var calls []any
	handler := &testHandler{
		handler: func(ctx context.Context, payload any) {
			calls = append(calls, payload)
		},
	}

	// Подписываемся на событие
	eventType := eventbus.EventType(1)
	bus.Subscribe(eventType, handler)

	// Добавляем события в очередь
	events.Enqueue(eventbus.NewEvent(eventType, "flush payload 1"))
	events.Enqueue(eventbus.NewEvent(eventType, "flush payload 2"))

	// Вызываем Flush
	bus.Flush(ctx, events)

	// Проверяем, что оба события были обработаны
	assert.Len(t, calls, 2)
	assert.Equal(t, "flush payload 1", calls[0])
	assert.Equal(t, "flush payload 2", calls[1])
}

func TestAsyncPublish(t *testing.T) {
	bus := eventbus.New()
	ctx := context.Background()
	wg := &sync.WaitGroup{}

	wg.Add(1)

	// Создаем обработчик
	var calls []any
	handler := &testHandler{
		handler: func(ctx context.Context, payload any) {
			defer wg.Done()
			calls = append(calls, payload)
		},
	}

	// Подписываемся на событие
	eventType := eventbus.EventType(1)
	bus.Subscribe(eventType, handler)

	// Публикуем асинхронное событие
	bus.Publish(ctx, eventbus.NewEvent(eventType, "async payload", eventbus.WithIsAsync(true)))

	// Даем время для обработки
	// Ждем немного, чтобы асинхронное событие было обработано
	wg.Wait()

	// Проверяем, что событие было обработано
	assert.Len(t, calls, 1)
	assert.Equal(t, "async payload", calls[0])
}

func TestUnsubscribe(t *testing.T) {
	bus := eventbus.New()
	ctx := context.Background()

	// Создаем обработчик
	var calls []any
	handler := &testHandler{
		handler: func(ctx context.Context, payload any) {
			calls = append(calls, payload)
		},
	}

	// Подписываемся на событие
	eventType := eventbus.EventType(1)
	handlerID, unsubscribe := bus.Subscribe(eventType, handler)

	assert.Equal(t, uint16(1), handlerID)

	// Публикуем событие
	bus.Publish(ctx, eventbus.NewEvent(eventType, "payload before unsubscribe"))

	// Убедимся, что обработчик был вызван
	assert.Len(t, calls, 1)

	// Отписываемся
	unsubscribe()

	// Публикуем еще одно событие
	bus.Publish(ctx, eventbus.NewEvent(eventType, "payload after unsubscribe"))

	// Убедимся, что обработчик не был вызван
	assert.Len(t, calls, 1)
}
