package main

import (
	"image"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func updateHoveredUIElement(plr *player) {
	plr.currentHoveredCard = nil
	plr.currentHoveredUIElement = UIEL_NONE

	v := rl.GetMousePosition()
	x, y := int32(v.X), int32(v.Y)
	mousePoint := image.Point{int(x), int(y)}

	if mousePoint.In(MulliganButton) {
		plr.currentHoveredUIElement = UIEL_MULLIGAN
		return
	}
	if mousePoint.In(EndTurnButton) {
		plr.currentHoveredUIElement = UIEL_ENDTURN
		return
	}

	if x >= CARDS_X && y >= CARDS_Y {
		index := (x - CARDS_X) / CARD_W
		if int(index) < plr.hand.Size() && index >= 0 {
			plr.currentHoveredUIElement = UIEL_HAND
			plr.currentHoveredCard = plr.hand[index]
			return
		}
	}
}

func handlePlayerClick(b *battlefield, plr *player) {
	if !rl.IsMouseButtonReleased(rl.MouseLeftButton) {
		return
	}
	if plr.currentHoveredUIElement == UIEL_ENDTURN {
		b.handlePlayerEndTurn(plr)
		return
	}
	if plr.currentHoveredUIElement == UIEL_MULLIGAN {
		b.handlePlayerMulligan(plr)
		return
	}
	if plr.currentHoveredCard != nil {
		b.playPlayerCard(plr, plr.currentHoveredCard)
		return
	}
}
