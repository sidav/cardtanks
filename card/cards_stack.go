package card

import "math/rand"

type CardsStack []*Card

func (s *CardsStack) Size() int {
	return len(*s)
}

func (s *CardsStack) Pop() *Card {
	// fmt.Printf("DEBUG: %d\n", len(*s))
	c := (*s)[0]
	*s = (*s)[1:]
	return c
}

func (s *CardsStack) TakeTopCardFromOtherStack(s2 *CardsStack) {
	c := s2.Pop()
	s.PushOnTop(c)
}

func (s *CardsStack) RemoveCard(c *Card) {
	for i := range *s {
		if (*s)[i] == c {
			*s = append((*s)[:i], (*s)[i+1:]...)
			return
		}
	}
	panic("No card " + c.Title + " in stack!")
}


func (s *CardsStack) PushOnTop(c *Card) {
	*s = append([]*Card{c}, *s...)
}

func (s *CardsStack) AddToBottom(c *Card) {
	*s = append(*s, c)
}


func (s CardsStack) Shuffle() {
	// Fisher–Yates shuffle
	for i := len(s) - 1; i > 0; i-- {
		exchInd := rand.Intn(i + 1)
		t := s[exchInd]
		s[exchInd] = s[i]
		s[i] = t
	}
}
