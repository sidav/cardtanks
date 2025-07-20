package card

type Card struct {
	Id                 CardId
	ActionsCost        int
	Title, description string
	ExhaustedOnUse     bool
}

func (c *Card) GetDescription() string {
	if c.ExhaustedOnUse {
		return c.description + "\nExhaust."
	}
	return c.description
}

