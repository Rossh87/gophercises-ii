package main

import (
	"fmt"
	"gophercises/deck/deck"
)

func main() {
	d := deck.New(deck.WithOmittedSuits(deck.Hearts, deck.Spades), deck.WithJokers(2), deck.WithOmittedValues(deck.Deuce))

	for _, card := range d {
		fmt.Printf("%s of %s\n", card.Value.Print(), card.Suit.Print())
	}
}
