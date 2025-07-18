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
	cardIds := [3]card.CardId{0, 0, 0}
	for i := range len(cardIds) {
		repeated := true
		for repeated {
			repeated = false
			generatedId := card.GetRandomCardId()
			if card.IsCardIdCommon(generatedId) && rand.Intn(2) == 0 {
				generatedId = card.GetRandomCardId()
			}
			cardIds[i] = generatedId
			for j := range len(cardIds) {
				if (i != j) && cardIds[i] == cardIds[j] {
					repeated = true
				}
			}
		}
	}
	for i, id := range cardIds {
		r.rewardCards[i] = card.CreateCardById(id)
	}
	return r
}

func (r *GameReward) GetDescription() string {
	return "Select 1 from 3 cards to add to your deck"
}
