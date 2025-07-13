package main

import (
	"fmt"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (r *renderer) renderGameState(g *game) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)

	switch g.state {
	case GS_BRIEFING:
		title, desc := g.bf.GetMissionDescription()
		r.drawTextInRect("Mission: "+title+
			" \n "+desc+
			" \n Reward: "+g.currentReward.GetDescription()+
			" \n Press ENTER to continue", 50, 50, 800, 25, rl.White)
	case GS_BATTLEFIELD:
		r.renderBattlefield(g.bf, g.plr)
	case GS_SELECT_REWARD:
		r.renderRewardScreen(g)
	case GS_GAMEOVER:
		r.drawTextInRect("GAME OVER \n press ESC to exit game", 50, 50, 800, 25, rl.Red)
	}

	rl.EndDrawing()
}

func (r *renderer) renderRewardScreen(g *game) {
	r.drawString(fmt.Sprintf("Select a card to add to your deck (you have %d cards)", g.plr.deck.Size()), 0, 0, 30, rl.White)
	for i, c := range g.currentReward.rewardCards {
		xPos := 50 + i*(CARD_W+CARD_W/4)
		r.drawString(strconv.Itoa(i+1), int32(xPos), 60, 33, rl.White)
		r.renderCardAt(c, xPos, 100, false)
	}
}
