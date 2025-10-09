package ws

import (
	"log"
	"net/http"
)

func ServeWS(h *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade error: %v", err)
		return
	}

	playerID := r.URL.Query().Get("playerId")
	gameID := r.URL.Query().Get("gameId")
	if playerID == "" || gameID == "" {
		_ = conn.Close()
		http.Error(w, "missing playerId/gameId", http.StatusBadRequest)
		return
	}

	c := NewClient(h, conn, playerID, gameID)
	h.register <- c

	go c.WritePump()
	go c.ReadPump()
}
