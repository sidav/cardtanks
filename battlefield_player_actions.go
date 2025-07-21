package main

import (
	"cardtanks/calc"
	"cardtanks/card"
	"math/rand"
)

func (b *battlefield) playPlayerCard(plr *player, c *card.Card) {
	pTank := b.playerTank
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
		b.state.intentVector.SetEqTo(b.playerTank.getFacing())
		b.state.actionsRemaining = 1

	case card.CARD_MOV2:
		b.state.switchTo(BS_PLAYER_MOVES)
		b.state.intentVector.SetEqTo(b.playerTank.getFacing())
		b.state.actionsRemaining = 2

	case card.CARD_MOV3:
		b.state.switchTo(BS_PLAYER_MOVES)
		b.state.intentVector.SetEqTo(b.playerTank.getFacing())
		b.state.actionsRemaining = 3

	case card.CARD_MOVLEFT:
		facingVector := calc.NewIntVector2d(pTank.getFacing())
		facingVector.RotateCCW()
		b.state.switchTo(BS_PLAYER_MOVES)
		b.state.intentVector = facingVector
		b.state.actionsRemaining = 1

	case card.CARD_MOVRIGHT:
		facingVector := calc.NewIntVector2d(pTank.getFacing())
		facingVector.RotateCW()
		b.state.switchTo(BS_PLAYER_MOVES)
		b.state.intentVector = facingVector
		b.state.actionsRemaining = 1

	case card.CARD_INSTANTSHOOT:
		b.state.switchTo(BS_PLAYER_SHOOTS_DURING_TURN)

	case card.CARD_DRAWCARD:
		if plr.hand.Size() >= plr.maxHandSize {
			return
		}
		for plr.hand.Size() < plr.maxHandSize {
			plr.drawCard()
		}
		plr.actionsSpentForTurn += c.ActionsCost
		plr.cardsPlayedThisTurn++
		if c.ExhaustedOnUse {
			plr.exhaustCard(c)
		} else {
			plr.discardCard(c)
		}
		b.handlePlayerEndTurn(plr)
		return

	case card.CARD_FRIENDLYTANK:
		x, y := b.playerTank.getCoordsFacingAt()
		if !(b.tileAt(x, y).is(TILE_FLOOR) &&
			b.getTankAt(x, y) == nil) {
			return
		}
		friend := createTank(TANK_PLAYER, TEAM_PLAYER, x, y)
		friend.faceRandomDirection()
		b.tanks = append(b.tanks, friend)
		b.state.pauseFor(300)

	case card.CARD_BUILD_WALLS_AROUND:
		x, y := pTank.getCoords()
		for i := x - 1; i <= x+1; i++ {
			for j := y - 1; j <= y+1; j++ {
				if b.areCoordsValid(i, j) && b.getTankAt(i, j) == nil {
					if b.tileAt(i, j).isOneOf(TILE_FLOOR, TILE_WATER, TILE_DAMAGED_WALL) {
						b.tileAt(i, j).spawnAs(TILE_WALL)
					}
				}
			}
		}
		b.state.pauseFor(300)

	case card.CARD_ROTATE_EVERYONE_RANDOMLY:
		b.playerTank.faceRandomDirection()
		b.playerTank.justSpawned = true
		for _, t := range b.tanks {
			t.faceRandomDirection()
			t.justSpawned = true
		}
		b.state.pauseFor(300)

	case card.CARD_ASSIGN_RANDOM_TEAMS:
		for _, t := range b.tanks {
			if t == b.enemyBossTank {
				continue
			}
			newTeam := TEAM_PLAYER
			for newTeam == TEAM_PLAYER || newTeam == TEAM_NONE {
				newTeam = rand.Intn(MAX_TEAM_CONST)
			}
			t.team = byte(newTeam)
			t.justSpawned = true
		}
		b.state.pauseFor(300)

	case card.CARD_SAFE_TELEPORT:
		v := b.selectRandomMapCoordsByAllowanceFunc(func(x, y int) bool {
			if !b.tileAt(x, y).isOneOf(TILE_FLOOR, TILE_ICE) {
				return false
			}
			for _, t := range b.tanks {
				if b.areTanksEnemies(b.playerTank, t) && b.lineOfFireExistsBetweenCoords(x, y, t.x, t.y) {
					return false
				}
			}
			return true
		})
		if v == nil {
			return
		}
		b.playerTank.justSpawned = true
		b.playerTank.x, b.playerTank.y = v.Unwrap()
		b.state.pauseFor(300)

	case card.CARD_UNEXHAUST_OTHER_CARDS:
		for plr.exhaustStack.Size() > 0 {
			plr.discard.TakeTopCardFromOtherStack(&plr.exhaustStack)
		}

	default:
		panic("Card '" + c.Title + "' not implemented")
	}

	plr.actionsSpentForTurn += c.ActionsCost
	plr.cardsPlayedThisTurn++
	if c.ExhaustedOnUse {
		plr.exhaustCard(c)
	} else {
		plr.discardCard(c)
	}
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
