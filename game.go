package main

import (
	"math/rand"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type gameStateId byte

const (
	GS_SELECT_BATTLE gameStateId = iota
	GS_BRIEFING
	GS_BATTLEFIELD
	GS_SELECT_REWARD
	GS_GAMEOVER
)

type game struct {
	state       gameStateId
	plr         *player
	bf          *battlefield
	r           *renderer
	wonMissions int

	currentReward *GameReward
}

func NewGame() *game {
	g := &game{
		state: GS_SELECT_REWARD,
		plr:   newPlayer(),
		r:     &renderer{},
		currentReward: GenerateNewReward(),
	}
	return g
}

func (g *game) Run() {
	for !rl.WindowShouldClose() {
		switch g.state {
		case GS_SELECT_BATTLE:
			g.selectBattle()
			g.state = GS_BRIEFING
		case GS_BRIEFING:
			k := rl.GetKeyPressed()
			if k == rl.KeyEnter {
				g.state = GS_BATTLEFIELD
				g.plr.drawInitialHand()
			}
		case GS_BATTLEFIELD:
			g.battlefieldLoop()
		case GS_SELECT_REWARD:
			g.SelectReward()
		case GS_GAMEOVER:
			// do nothing
		}
		g.r.renderGameState(g)
	}
}

func (g *game) selectBattle() {
	mission := rand.Intn(int(BFM_MISSIONS_COUNT))
	switch battlefieldMissionId(mission) {
	case BFM_SKIRMISH:
		g.bf = createBattlefieldSkirmish(0, 3+g.wonMissions/2, 5+g.wonMissions)
	case BFM_CAPTURE_FLAGS:
		g.bf = createBattlefieldCaptureFlags(0, 3+g.wonMissions/2)
	case BFM_DESTROY_EAGLES:
		g.bf = createBattlefieldDestroyEagles(0, 3+g.wonMissions/2)
	default:
		panic("No mission implementation")
	}
	// Generate reward
	g.currentReward = GenerateNewReward()
}

func (g *game) battlefieldLoop() {
	g.bf.actOnState(g.plr)
	updateHoveredUIElement(g.plr)
	if g.bf.state.awaitsPlayerInput() {
		handlePlayerClick(g.bf, g.plr)
	}
	if g.bf.IsMissionLost() {
		g.state = GS_GAMEOVER
		return
	}
	if g.bf.IsMissionWon() {
		g.plr.returnAllCardsToDeck()
		g.wonMissions++
		g.state = GS_SELECT_REWARD
		return
	}
}

func (g *game) SelectReward() {
	k := rl.GetKeyPressed()
	for i := range g.currentReward.rewardCards {
		if g.IsKeyCodeEqualToString(k, strconv.Itoa(i+1), false) {
			g.plr.deck.PushOnTop(g.currentReward.rewardCards[i])
			g.state = GS_SELECT_BATTLE
		}
	}
}

func (g *game) IsKeyCodeEqualToString(keyCode int32, keyString string, withShift bool) bool {
	shiftPressed := rl.IsKeyDown(rl.KeyLeftShift) || rl.IsKeyDown(rl.KeyRightShift)
	return int32(keyString[0])-keyCode == 0 && withShift == shiftPressed
}
