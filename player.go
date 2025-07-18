package main

import (
	"cardtanks/card"
	"sort"
)

type UIElementCode byte

const (
	UIEL_NONE UIElementCode = iota
	UIEL_HAND
	UIEL_ENDTURN
	UIEL_MULLIGAN
)

type player struct {
	maxHandSize  int
	deck         card.CardsStack
	hand         card.CardsStack
	discard      card.CardsStack
	exhaustStack card.CardsStack

	currentHoveredUIElement UIElementCode
	currentHoveredCard      *card.Card

	didMulliganThisTurn bool
	cardsPlayedThisTurn int
	actionsSpentForTurn int
}

func newPlayer() *player {
	plr := &player{
		maxHandSize: 5,
	}
	plr.deck.PushOnTop(card.CreateCardById(card.CARD_ROTLEFT))
	plr.deck.PushOnTop(card.CreateCardById(card.CARD_ROTLEFT))
	plr.deck.PushOnTop(card.CreateCardById(card.CARD_ROTRIGHT))
	plr.deck.PushOnTop(card.CreateCardById(card.CARD_ROTRIGHT))
	plr.deck.PushOnTop(card.CreateCardById(card.CARD_MOV1))
	plr.deck.PushOnTop(card.CreateCardById(card.CARD_MOV2))
	plr.deck.PushOnTop(card.CreateCardById(card.CARD_MOV3))

	return plr
}

func (p *player) sortHand() {
	sort.Slice(p.hand, func(i, j int) bool { return (p.hand)[i].Id < (p.hand)[j].Id })
}

func (p *player) returnAllCardsToDeck() {
	for p.hand.Size() > 0 {
		p.deck.TakeTopCardFromOtherStack(&p.hand)
	}
	for p.discard.Size() > 0 {
		p.deck.TakeTopCardFromOtherStack(&p.discard)
	}
	for p.exhaustStack.Size() > 0 {
		p.deck.TakeTopCardFromOtherStack(&p.exhaustStack)
	}
}

func (p *player) drawInitialHand() {
	p.cardsPlayedThisTurn = 0
	p.actionsSpentForTurn = 0
	p.deck.Shuffle()
	for p.hand.Size() < p.maxHandSize {
		c := p.deck.Pop()
		p.hand.PushOnTop(c)
	}
	p.sortHand()
}

func (p *player) drawCard() {
	if p.deck.Size() == 0 {
		for p.discard.Size() > 0 {
			p.deck.TakeTopCardFromOtherStack(&p.discard)
		}
		p.deck.Shuffle()
	}
	if p.deck.Size() > 0 {
		p.hand.TakeTopCardFromOtherStack(&p.deck)
	}
	p.sortHand()
}

func (p *player) discardCard(c *card.Card) {
	p.hand.RemoveCard(c)
	p.discard.PushOnTop(c)
}

func (p *player) exhaustCard(c *card.Card) {
	p.hand.RemoveCard(c)
	p.exhaustStack.PushOnTop(c)
}

func (p *player) discardHand() {
	for p.hand.Size() > 0 {
		p.discardCard(p.hand[0])
	}
}
