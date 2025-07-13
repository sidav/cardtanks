package main

import (
	"cardtanks/card"
	"math/rand"
)

type GameReward struct {
	typeId      int
	rewardCards [3]*card.Card
}

func GenerateNewReward() *GameReward {
	r := &GameReward{}
	// Generate card rewards
	cardIds := [3]int{0, 0, 0}
	for i := range len(cardIds) {
		repeated := true
		for repeated {
			repeated = false
			cardIds[i] = rand.Intn(int(card.TOTAL_CARDS))
			for j := range len(cardIds) {
				if (i != j) && cardIds[i] == cardIds[j] {
					repeated = true
				}
			}
		}
	}
	for i, id := range cardIds {
		r.rewardCards[i] = card.CreateCardById(card.CardId(id))
	}
	return r
}

func (r *GameReward) GetDescription() string {
	return "Select 1 from 3 cards to add to your deck"
}

