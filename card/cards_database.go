package card

import "math/rand"

type CardId int

const (
	// Common
	CARD_ROTLEFT CardId = iota
	CARD_ROTRIGHT
	CARD_ROTAROUND
	CARD_MOV1
	CARD_MOV2
	CARD_MOV3
	// Uncommon
	CARD_MOVLEFT
	CARD_MOVRIGHT
	CARD_INSTANTSHOOT
	CARD_DRAWCARD_AND_REPAIR
	CARD_FRIENDLYTANK
	CARD_BUILD_WALLS_AROUND
	CARD_PUSH_AND_ROTATE
	CARD_ASSIGN_RANDOM_TEAMS
	CARD_SAFE_TELEPORT
	CARD_TELEPORT_BEHIND
	CARD_GROW_FOREST
	CARD_REMOVE_WALL
	CARD_UNEXHAUST_OTHER_CARDS
	TOTAL_CARDS
)

func IsCardIdCommon(id CardId) bool {
	return id <= CARD_MOV3
}

func GetRandomCardId() CardId {
	return CardId(rand.Intn(int(TOTAL_CARDS)))
}

func CreateCardById(id CardId) *Card {
	var c *Card
	switch id {
	case CARD_ROTLEFT:
		c = &Card{
			Title: "Turn left",
		}
	case CARD_ROTRIGHT:
		c = &Card{
			Title: "Turn right",
		}
	case CARD_ROTAROUND:
		c = &Card{
			Title: "Turn around",
		}
	case CARD_MOV1:
		c = &Card{
			ActionsCost: 1,
			Title:       "Move 1",
		}
	case CARD_MOV2:
		c = &Card{
			ActionsCost: 2,
			Title:       "Move 2",
		}
	case CARD_MOV3:
		c = &Card{
			ActionsCost: 3,
			Title:       "Move 3",
		}
	case CARD_MOVLEFT:
		c = &Card{
			ActionsCost: 0,
			Title:       "Left sidestep",
			description: "Move left without turning",
		}
	case CARD_MOVRIGHT:
		c = &Card{
			ActionsCost: 0,
			Title:       "Right sidestep",
			description: "Move right without turning",
		}
	case CARD_INSTANTSHOOT:
		c = &Card{
			ActionsCost:    2,
			Title:          "Instant fire!",
			description:    "Shoot right now without ending your turn.",
			ExhaustedOnUse: true,
		}
	case CARD_DRAWCARD_AND_REPAIR:
		c = &Card{
			ActionsCost:    3,
			Title:          "Quick repair",
			description:    "If damaged, repair 1 health, draw a card and end your turn.",
			ExhaustedOnUse: true,
		}
	case CARD_FRIENDLYTANK:
		c = &Card{
			ActionsCost:    3,
			Title:          "Call for help",
			description:    "Summon a light friendly tank on empty square before you.",
			ExhaustedOnUse: true,
		}
	case CARD_BUILD_WALLS_AROUND:
		c = &Card{
			ActionsCost:    3,
			Title:          "Bunker",
			description:    "Build/repair walls around you, except the tile you're looking at",
			ExhaustedOnUse: true,
		}
	case CARD_PUSH_AND_ROTATE:
		c = &Card{
			ActionsCost:    0,
			Title:          "Telekinesis",
			description:    "Push and randomly rotate the tank you're looking at",
			ExhaustedOnUse: true,
		}
	case CARD_ASSIGN_RANDOM_TEAMS:
		c = &Card{
			ActionsCost:    3,
			Title:          "Mayhem",
			description:    "Assign random teams to most tanks.",
			ExhaustedOnUse: true,
		}
	case CARD_UNEXHAUST_OTHER_CARDS:
		c = &Card{
			ActionsCost:    2,
			Title:          "Second wind",
			description:    "Return all exhausted cards to your discard pile.",
			ExhaustedOnUse: true,
		}
	case CARD_SAFE_TELEPORT:
		c = &Card{
			ActionsCost:    2,
			Title:          "Evacuate",
			description:    "Try to teleport to a safe place.",
			ExhaustedOnUse: true,
		}
	case CARD_TELEPORT_BEHIND:
		c = &Card{
			ActionsCost:    0,
			Title:          "Sneak attack",
			description:    "Teleport behind the tank you're looking at",
			ExhaustedOnUse: true,
		}
	case CARD_GROW_FOREST:
		c = &Card{
			ActionsCost:    2,
			Title:          "Foliage camo",
			description:    "Grow a forest on your position",
			ExhaustedOnUse: true,
		}
	case CARD_REMOVE_WALL:
		c = &Card{
			ActionsCost:    0,
			Title:          "Unbuild",
			description:    "Remove any wall you're looking at (works even with armor)",
			ExhaustedOnUse: true,
		}
	default:
		c = &Card{
			Title:       "Unimplemented card",
			description: "If you see this, there is a bug",
		}
	}

	c.Id = id
	return c
}
