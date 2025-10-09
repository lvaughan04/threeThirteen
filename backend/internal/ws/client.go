package ws

import (
	"backend/internal/events/ui"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024
)

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte //outgoing messages
	playerID string
	gameID   string
}

func NewClient(h *Hub, conn *websocket.Conn, playerID, gameID string) *Client {
	return &Client{
		hub:      h,
		conn:     conn,
		send:     make(chan []byte, 64),
		playerID: playerID,
		gameID:   gameID,
	}
}

func (c *Client) ReadPump() {
	defer func() {
		// on any exit: unregister and let hub close send (which ends writePump)
		c.hub.unregister <- c
		_ = c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// Parse incoming message and convert to Event
		var rawMsg map[string]interface{}
		if err := json.Unmarshal(message, &rawMsg); err != nil {
			log.Printf("Failed to parse message: %v", err)
			continue
		}

		// Route to EventBus based on message type
		c.routeMessage(rawMsg)
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				if _, err := w.Write([]byte{'\n'}); err != nil {
					_ = w.Close()
					return
				}
				if _, err := w.Write(<-c.send); err != nil {
					_ = w.Close()
					return
				}
			}
			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) routeMessage(msg map[string]interface{}) {
	msgType, ok := msg["type"].(string)
	if !ok {
		return
	}

	switch msgType {
	case ui.TT_PLAYER_DREW_CARD:
		event := ui.PlayerDrewCardIntent{
			GameID:   msg["game_id"].(string),
			PlayerID: msg["player_id"].(string),
			Source:   msg["source"].(string),
		}
		c.hub.eventBus.Publish(event)

	case ui.TT_PLAYER_DISCARDED_CARD:
		event := ui.PlayerDiscardedCardIntent{
			GameID:   msg["game_id"].(string),
			PlayerID: msg["player_id"].(string),
			Card:     msg["card"],
		}
		c.hub.eventBus.Publish(event)
	}
}
