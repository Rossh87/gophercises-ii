package blackjack

import "fmt"

type CardValue int

const (
	Ace CardValue = iota
	Deuce
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Joker
)

func (s CardValue) Print() string {
	return []string{"Ace", "Deuce", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Jack", "Queen", "King", "Joker"}[s]
}

func (s CardValue) Value() int {
	return []int{0, 2, 3, 4, 5, 6, 7, 8, 9, 10, 10, 10, 10, 10}[s]
}

type Suit int

const (
	Spades Suit = iota
	Diamonds
	Clubs
	Hearts
	Red
	Black
)

func (s Suit) Print() string {
	return []string{"Spades", "Diamonds", "Clubs", "Hearts", "Red", "Black"}[s]
}

type Card struct {
	Suit  Suit
	Value CardValue
}

func (c Card) Print() string {
	return fmt.Sprintf("%s of %s\n", c.Value.Print(), c.Suit.Print())
}
