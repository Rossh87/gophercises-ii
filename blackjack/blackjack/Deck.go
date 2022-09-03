package blackjack

import (
	"math/rand"
	"time"
)

type Deck struct {
	head  int
	cards []Card
	empty bool
}

func (d *Deck) Empty() bool {
	return d.empty
}

func (d *Deck) Deal() (*Card, error) {
	if d.head == len(d.cards) {
		panic("Attempting to deal from an empty deck")
	}

	c := d.cards[d.head]

	d.head++

	if d.head == len(d.cards) {
		d.empty = true
	}

	return &c, nil
}

func (d *Deck) Shuffle() {
	rand.Seed(time.Now().Unix())

	rand.Shuffle(len(d.cards)-d.head, func(i, j int) {
		offseti := i + d.head
		offfsetj := j + d.head
		d.cards[offseti], d.cards[offfsetj] = d.cards[offfsetj], d.cards[offseti]
	})
}

type deckConf struct {
	sort       func(d []Card) []Card
	customSort bool
	jokers     int
	omitSuit   []Suit
	omitVal    []CardValue
	decks      int
}

type deckOpt func(conf *deckConf)

func WithSort(f func(d []Card) []Card) deckOpt {
	return func(conf *deckConf) {
		conf.sort = f
		conf.customSort = true
	}
}

func WithJokers(jokers int) deckOpt {
	return func(conf *deckConf) {
		conf.jokers = jokers
	}
}

func WithOmittedSuits(omitted ...Suit) deckOpt {
	return func(conf *deckConf) {
		conf.omitSuit = omitted
	}
}

func WithOmittedValues(omitted ...CardValue) deckOpt {
	return func(conf *deckConf) {
		conf.omitVal = omitted
	}
}

func WithDecks(decks int) deckOpt {
	return func(conf *deckConf) {
		conf.decks = decks
	}
}

func NewDeck(opts ...deckOpt) Deck {
	conf := deckConf{
		func(d []Card) []Card {
			return d
		},
		false,
		0,
		[]Suit{},
		[]CardValue{},
		1,
	}

	for _, opt := range opts {
		opt(&conf)
	}

	omittedSuits := make(map[Suit]struct{})
	omittedValues := make(map[CardValue]struct{})

	for _, s := range conf.omitSuit {
		omittedSuits[s] = struct{}{}
	}

	for _, v := range conf.omitVal {
		omittedValues[v] = struct{}{}
	}

	var combinedCards []Card

	for dc := 0; dc < conf.decks; dc++ {
		cards := []Card{}

		for s := 0; s < 4; s++ {
			for v := 0; v < 13; v++ {
				suit := Suit(s)
				value := CardValue(v)

				if _, ok := omittedSuits[suit]; ok {
					continue
				}

				if _, ok := omittedValues[value]; ok {
					continue
				}

				cards = append(cards, Card{suit, value})
			}
		}

		for j := 0; j < conf.jokers; j++ {
			var suit Suit

			if j&1 == 1 {
				suit = Red
			} else {
				suit = Black
			}

			cards = append(cards, Card{suit, Joker})
		}

		combinedCards = append(combinedCards, cards...)
	}

	if conf.customSort {
		combinedCards = conf.sort(combinedCards)
	}

	deck := Deck{
		head:  0,
		cards: combinedCards,
		empty: false,
	}

	return deck
}
