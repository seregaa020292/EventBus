package eventbus_test

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"sanbox/eventbus"
)

type testHandler struct {
	calls []eventbus.Event
	mu    sync.Mutex
}

func (h *testHandler) Handle(ctx context.Context, event eventbus.Event) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.calls = append(h.calls, event)
}

type testEvent string

func (e testEvent) Name() eventbus.EventName { return 1 }

func TestEventBus(t *testing.T) {
	bus := eventbus.New()
	ctx := context.Background()

	// Создаем обработчик
	handler := &testHandler{}

	// Подписываемся на событие
	event := testEvent("test payload")
	_, unsubscribe := bus.Subscribe(event.Name(), handler)

	// Публикуем событие
	bus.Publish(ctx, event)

	// Убедимся, что обработчик был вызван
	assert.Len(t, handler.calls, 1)
	assert.Equal(t, testEvent("test payload"), handler.calls[0])

	// Отписываемся и проверяем, что обработчик больше не вызывается
	unsubscribe()
	bus.Publish(ctx, testEvent("another payload"))

	// Убедимся, что обработчик не был вызван
	assert.Len(t, handler.calls, 1)
}

func TestFlush(t *testing.T) {
	bus := eventbus.New()
	ctx := context.Background()

	// Создаем обработчик
	handler := &testHandler{}

	// Подписываемся на событие
	event := testEvent("test payload")
	bus.Subscribe(event.Name(), handler)

	// Добавляем события в очередь
	events := &eventbus.Events{}
	events.Enqueue(event)
	events.Enqueue(event)

	// Вызываем Flush
	bus.Flush(ctx, events)

	// Проверяем, что оба события были обработаны
	assert.Len(t, handler.calls, 2)
	assert.Equal(t, testEvent("test payload"), handler.calls[0])
	assert.Equal(t, testEvent("test payload"), handler.calls[1])
}

func TestAsyncPublish(t *testing.T) {
	bus := eventbus.New()
	ctx := context.Background()

	// Создаем обработчик
	handler := &testHandler{}

	// Подписываемся на событие
	event := testEvent("async payload")
	bus.Subscribe(event.Name(), handler, eventbus.WithHandlerIsAsync(true))

	// Публикуем асинхронное событие
	bus.Publish(ctx, event)

	// Даем время для обработки
	bus.Wait()

	// Проверяем, что событие было обработано
	assert.Len(t, handler.calls, 1)
	assert.Equal(t, testEvent("async payload"), handler.calls[0])
}

func TestUnsubscribe(t *testing.T) {
	bus := eventbus.New()
	ctx := context.Background()

	// Создаем обработчик
	handler := &testHandler{}

	// Подписываемся на событие
	event := testEvent("payload unsubscribe")
	handlerID, unsubscribe := bus.Subscribe(event.Name(), handler)

	// Публикуем событие
	bus.Publish(ctx, event)

	// Убедимся, что обработчик был вызван
	assert.Equal(t, uint16(1), handlerID)
	assert.Len(t, handler.calls, 1)

	// Отписываемся
	unsubscribe()

	// Публикуем еще одно событие
	bus.Publish(ctx, event)

	// Убедимся, что обработчик не был вызван
	assert.Len(t, handler.calls, 1)
}
