package hub

import (
	"backend/internal/events/game"
	"sync"
)

type EventBus struct {
	subscribers map[string][]func(game.GameEvent)
	mu          sync.RWMutex
}

func (eb *EventBus) Subscribe(eventType string, handler func(game.GameEvent)) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.subscribers[eventType] = append(eb.subscribers[eventType], handler)
}

func (eb *EventBus) Publish(event game.GameEvent) {    
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	if handlers, found := eb.subscribers[event.Type]; found {
		for _, handler := range handlers {
			go handler(event)
		}
	}
}
