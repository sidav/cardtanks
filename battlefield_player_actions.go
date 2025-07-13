package main

import (
	"cardtanks/calc"
	"cardtanks/card"
	"math/rand"
)

func (b *battlefield) playPlayerCard(plr *player, c *card.Card) {
	if !b.canCardBePlayed(plr, c) {
		return
	}
	switch c.Id {
	case card.CARD_ROTLEFT:
		b.playerTank.rotateLeft()

	case card.CARD_ROTRIGHT:
		b.playerTank.rotateRight()

	case card.CARD_ROTAROUND:
		b.playerTank.rotateLeft()
		b.playerTank.rotateLeft()
	case card.CARD_WAIT:
		// Nothing
	case card.CARD_MOV1:
		b.state.switchTo(BS_PLAYER_MOVES)
		b.state.actionsRemaining = 1

	case card.CARD_MOV2:
		b.state.switchTo(BS_PLAYER_MOVES)
		b.state.actionsRemaining = 2

	case card.CARD_MOV3:
		b.state.switchTo(BS_PLAYER_MOVES)
		b.state.actionsRemaining = 3

	case card.CARD_DRAWCARD:
		if plr.hand.Size() >= plr.maxHandSize {
			return
		}
		plr.drawCard()
		plr.actionsSpentForTurn += c.ActionsCost
		plr.cardsPlayedThisTurn++
		plr.discardCard(c)
		b.handlePlayerEndTurn(plr)
		return

	case card.CARD_FRIENDLYTANK:
		v := b.selectRandomMapCoordsByAllowanceFunc(func(x, y int) bool {
			return b.tileAt(x, y).is(TILE_FLOOR) &&
				calc.ApproxDistanceInt(x, y, b.playerTank.x, b.playerTank.y) == 1 &&
				b.getTankAt(x, y) == nil
		})
		if v == nil {
			return
		}
		friend := createTank(TANK1, TEAM_PLAYER, v.X, v.Y)
		b.tanks = append(b.tanks, friend)
		b.state.switchTo(BS_TEMP_PAUSE)

	default:
		panic("Card '" + c.Title + "' not implemented")
	}

	plr.actionsSpentForTurn += c.ActionsCost
	plr.cardsPlayedThisTurn++
	plr.discardCard(c)
}

func (b *battlefield) canCardBePlayed(plr *player, c *card.Card) bool {
	if plr.actionsSpentForTurn > 0 && c.ActionsCost > 0 {
		return false
	}
	return true
}

func (b *battlefield) canPlayerMulligan(plr *player) bool {
	return plr.cardsPlayedThisTurn == 0 && !plr.didMulliganThisTurn
}

func (b *battlefield) handlePlayerMulligan(plr *player) {
	if !b.canPlayerMulligan(plr) {
		return
	}
	cardsBefore := plr.hand.Size()
	plr.discardHand()
	for range cardsBefore {
		plr.drawCard()
	}
	plr.didMulliganThisTurn = true
}

func (b *battlefield) handlePlayerEndTurn(plr *player) {
	cardsToDraw := plr.cardsPlayedThisTurn + plr.hand.Size()
	plr.cardsPlayedThisTurn = 0
	plr.discardHand()
	for cardsToDraw > 0 {
		plr.drawCard()
		cardsToDraw--
	}
	plr.didMulliganThisTurn = false
	b.state.switchTo(BS_PLAYER_ENDED_TURN)
}

const CARDS_DROPPED_ON_HIT = 2

func (b *battlefield) handlePlayerBeingHit(plr *player) {
	if plr.hand.Size() > CARDS_DROPPED_ON_HIT {
		for range CARDS_DROPPED_ON_HIT {
			ind := rand.Intn(plr.hand.Size())
			crd := plr.hand[ind]
			plr.discardCard(crd)
		}
	} else {
		// game over should be here
		b.playerTank.health = 0
	}
}
