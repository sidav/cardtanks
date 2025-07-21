package main

import (
	"cardtanks/card"
	"fmt"
	"image"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	CARD_W              = 175
	CARD_H              = 220
	CARDS_X             = 175
	CARDS_Y             = WINDOW_H - 4*CARD_H/5
	HOVERED_CARD_OFFSET = CARD_H / 5
)

var EndTurnButton = image.Rect(WINDOW_W-150, WINDOW_H-100, WINDOW_W, WINDOW_H)
var MulliganButton = image.Rect(0, WINDOW_H-100, 150, WINDOW_H)

func (r *renderer) renderBattlefieldInfo() {
	bfW := len(r.bf.tiles) * TILE_SIZE_PIXELS
	r.drawString(fmt.Sprintf("Current phase: %s", r.bf.state.currentStateName()), int32(bfW+4), 4, 24, rl.White)
	r.drawString(fmt.Sprintf("Health: %d/%d", r.bf.playerTank.health, r.bf.playerTank.GetMaxHealth()), int32(bfW+4), 30, 24, rl.White)
	r.drawString(r.bf.GetMissionProgressString(), int32(bfW+4), 60, 24, rl.White)
}

func (r *renderer) renderButtons() {
	// Deck/mulligan button
	color := rl.White
	if r.bf.canPlayerMulligan(r.plr) && r.plr.currentHoveredUIElement == UIEL_MULLIGAN {
		color = rl.Yellow
	}
	rl.DrawRectangleLinesEx(r.imgRectangeToRlRectangle(MulliganButton), 2, color)
	r.drawString(fmt.Sprintf("Deck (%d)", r.plr.deck.Size()), int32(MulliganButton.Min.X+2), int32(MulliganButton.Min.Y+2), 18, color)
	if r.bf.canPlayerMulligan(r.plr) {
		r.drawString("Press to mulligan", int32(MulliganButton.Min.X+2), int32(MulliganButton.Min.Y+22), 16, color)
	}

	// End turn button
	color = rl.White
	if r.plr.currentHoveredUIElement == UIEL_ENDTURN {
		color = rl.Yellow
	}
	rl.DrawRectangleLinesEx(r.imgRectangeToRlRectangle(EndTurnButton), 2, color)
	r.drawString("End turn", int32(EndTurnButton.Min.X), int32(EndTurnButton.Min.Y), 18, color)
	r.drawString(fmt.Sprintf("(%d actions)", r.plr.actionsSpentForTurn), int32(EndTurnButton.Min.X), int32(EndTurnButton.Min.Y+22), 18, color)
}

func (r *renderer) renderCardsInHand() {
	for i, c := range r.plr.hand {
		var yOffset int
		if r.plr.currentHoveredCard == c {
			yOffset = HOVERED_CARD_OFFSET
		}
		r.renderCardAt(c, CARDS_X+CARD_W*i+5, CARDS_Y-yOffset, true)
	}
}

func (r *renderer) renderCardAt(c *card.Card, x, y int, considerBfPlayability bool) {
	color := rl.LightGray
	if c.ActionsCost > 0 {
		color = rl.Orange
	}
	if considerBfPlayability && !r.bf.canCardBePlayed(r.plr, c) {
		color = rl.DarkGray
	}
	rl.DrawRectangleLinesEx(
		rl.Rectangle{
			X:      float32(x),
			Y:      float32(y),
			Width:  float32(CARD_W),
			Height: float32(CARD_H),
		},
		2, // thickness
		color,
	)
	r.drawString(fmt.Sprintf("(%d) %s", c.ActionsCost, c.Title),
		int32(x), int32(y), 18, color)
	r.drawTextInRect(c.GetDescription(), float32(x+1), float32(y+22), CARD_W, 16, color)
}

func (r *renderer) drawString(text string, x, y, size int32, color rl.Color) {
	rl.DrawTextEx(defaultFont, text, rl.Vector2{float32(x + 2), float32(y + 2)}, float32(size), 0, color)
}

func (r *renderer) drawTextInRect(text string, x, y, w, fontSize float32, color rl.Color) {
	spaceWidth := fontSize / 2
	lineMargin := fontSize / 5
	var currX float32 = 0
	var currY float32 = y
	newlineSplittedText := strings.Split(text, "\n")
	for _, line := range newlineSplittedText {
		spaceSplittedText := strings.Split(line, " ")
		for _, word := range spaceSplittedText {
			measure := rl.MeasureTextEx(defaultFont, word, fontSize, 0)
			wordWidth := measure.X
			if currX+wordWidth > w {
				currX = 0.0
				currY += fontSize + lineMargin
			}
			rl.DrawTextEx(defaultFont, word, rl.Vector2{currX + x, currY}, fontSize, 0, color)
			currX += wordWidth + spaceWidth
		}
		currX = 0.0
		currY += fontSize + lineMargin
	}
}

func (r *renderer) imgRectangeToRlRectangle(rect image.Rectangle) rl.Rectangle {
	return rl.Rectangle{
		X:      float32(rect.Min.X),
		Y:      float32(rect.Min.Y),
		Width:  float32(rect.Dx()),
		Height: float32(rect.Dy()),
	}
}
