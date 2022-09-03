package blackjack

type Player struct {
	id     int
	hand   []Card
	bust   bool
	stands bool
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func (p Player) PrintHand() string {
	return printhand(p.hand, nil)
}

func (p Player) Total() int {
	return total(p.hand)
}
