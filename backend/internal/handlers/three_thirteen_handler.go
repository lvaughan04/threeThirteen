package handlers

import (
	"backend/internal/events/interfaces"
	"backend/internal/events/ui"
	"backend/internal/events/game"
	"backend/internal/hub"
	"backend/internal/models"
	"backend/internal/types"
)

type ThreeThirteenHandler struct {
	eventBus   *hub.EventBus
	activeGames map[string]*models.ThreeThirteenGame
}

func (h *ThreeThirteenHandler) handlePlayerDrewCard(event ui.PlayerDrewCardIntent) {
	game, ok := h.activeGames[event.GetGameID()];
	if !ok {
		return
	}
	player , ok := game.Players[event.PlayerID]
	if !ok {
		return
	}
	var drawnCard types.Card
	if event.Source == "deck" {
		game.DrawFromDeck(player)
	} else if event.Source == "discard" {
		game.DrawFromDiscard(player)
	}

	h.eventBus.Publish(game_events.GameEvent{
		Type:   ui.TT_CARD_DRAWN,
		GameID: event.GameID,
		Data: ui.TTCardDrawnData{
			PlayerID: event.PlayerID,
			Source:   event.Source,
			Card:     
		},
	})
		
}

func (h *ThreeThirteenHandler) handlePlayerDiscardedCard(event ui.PlayerDiscardedCardIntent) {

}

func (h *ThreeThirteenHandler) handlePlayerWentOut(event ui.PlayerWentOutIntent) {

}

func (h *ThreeThirteenHandler) Register() {
	h.eventBus.Subscribe((ui.TT_PLAYER_DREW_CARD), func(e interfaces.Event) {
		h.handlePlayerDrewCard(e.(ui.PlayerDrewCardIntent))
	})
	h.eventBus.Subscribe((ui.TT_PLAYER_DISCARDED_CARD), func(e interfaces.Event) {
		h.handlePlayerDiscardedCard(e.(ui.PlayerDiscardedCardIntent))
	})
	h.eventBus.Subscribe((ui.TT_PLAYER_WENT_OUT), func(e interfaces.Event) {
		h.handlePlayerWentOut(e.(ui.PlayerWentOutIntent))
	})
}

func NewThreeThirteenHandler(eventBus *hub.EventBus) *ThreeThirteenHandler {
	handler := &ThreeThirteenHandler{
		eventBus:   eventBus,
		activeGames: make(map[string]*models.ThreeThirteenGame),
	}
	handler.Register()
	return handler
}



