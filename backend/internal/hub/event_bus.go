package hub

import (
	"backend/internal/events/interfaces"
	"sync"
)

type EventBus struct {
	subscribers map[string][]func(interfaces.Event)
	mu          sync.RWMutex
}

func (eb *EventBus) Subscribe(eventType string, handler func(interfaces.Event)) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.subscribers[eventType] = append(eb.subscribers[eventType], handler)
}

func (eb *EventBus) Publish(event interfaces.Event) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	if handlers, found := eb.subscribers[event.GetType()]; found {
		for _, handler := range handlers {
			go handler(event)
		}
	}
}
