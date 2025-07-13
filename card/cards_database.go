package card

type CardId int

const (
	CARD_ROTLEFT CardId = iota
	CARD_ROTRIGHT
	CARD_ROTAROUND
	CARD_WAIT
	CARD_MOV1
	CARD_MOV2
	CARD_MOV3
	CARD_DRAWCARD
	CARD_FRIENDLYTANK
	TOTAL_CARDS
)

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
	case CARD_WAIT:
		c = &Card{
			ActionsCost: 2,
			Title:       "Wait a while",
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
	case CARD_DRAWCARD:
		c = &Card{
			ActionsCost: 3,
			Title:       "Quick repair",
			Description: "If damaged, draw a card and end your turn.",
		}
	case CARD_FRIENDLYTANK:
		c = &Card{
			ActionsCost: 3,
			Title:       "Call for help",
			Description: "Summon a light friendly tank on empty square nearby.",
		}
	default:
		c = &Card{
			Title: "Unimplemented card",
		}
	}

	c.Id = id
	return c
}
