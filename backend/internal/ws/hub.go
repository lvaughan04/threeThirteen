package ws

import (
	"backend/internal/events/interfaces"
	"backend/internal/hub"
	"encoding/json"
	"log"
	"sync"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	eventBus   *hub.EventBus
	mu         sync.RWMutex
}

func NewHub(eventBus *hub.EventBus) *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		eventBus:   eventBus,
	}
}

// Run starts the hub - MUST run in goroutine
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Client registered. Total: %d", len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			log.Printf("Client unregistered. Total: %d", len(h.clients))

		case message := <-h.broadcast:
			var toRemove []*Client

			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					toRemove = append(toRemove, client)
				}
			}
			h.mu.RUnlock()
			if len(toRemove) > 0 {
				h.mu.Lock()
				for _, c := range toRemove {
					if _, ok := h.clients[c]; ok {
						delete(h.clients, c)
						close(c.send)
					}
				}
				h.mu.Unlock()
				log.Printf("Removed %d slow/dead clients", len(toRemove))
			}
		}
	}
}

func (h *Hub) BroadcastToGame(gameID string, message []byte) {
	var toRemove []*Client

	h.mu.RLock()
	for client := range h.clients {
		if client.gameID != gameID {
			continue
		}
		select {
		case client.send <- message:
			// Successfully sent
		default:
			// Channel full
			toRemove = append(toRemove, client)
		}
	}
	h.mu.RUnlock()

	// Remove dead clients with write lock
	if len(toRemove) > 0 {
		h.mu.Lock()
		for _, c := range toRemove {
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.send)
			}
		}
		h.mu.Unlock()
	}
}

func (h *Hub) SubscribeToGameEvents() {
	h.eventBus.Subscribe("TT_CARD_DRAWN", func(e interfaces.Event) {
		data, _ := json.Marshal(e)
		h.BroadcastToGame(e.GetGameID(), data)
	})

	h.eventBus.Subscribe("TT_CARD_DISCARDED", func(e interfaces.Event) {
		data, _ := json.Marshal(e)
		h.BroadcastToGame(e.GetGameID(), data)
	})
}
