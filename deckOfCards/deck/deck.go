package deck

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

type deckConf struct {
	sort     func(d []Card) []Card
	shuffle  bool
	jokers   int
	omitSuit []Suit
	omitVal  []CardValue
	decks    int
}

type deckOpt func(conf *deckConf)

func WithSort(f func(d []Card) []Card) deckOpt {
	return func(conf *deckConf) {
		conf.sort = f
	}
}

func WithShuffle(shuf bool) deckOpt {
	return func(conf *deckConf) {
		conf.shuffle = shuf
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

func New(opts ...deckOpt) []Card {
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

	var decks []Card

	for dc := 0; dc < conf.decks; dc++ {
		deck := []Card{}
		omitted := 0

		for s := 0; s < 4; s++ {
			for v := 0; v < 13; v++ {
				suit := Suit(s)
				value := CardValue(v)

				if _, ok := omittedSuits[suit]; ok {
					omitted++
					continue
				}

				if _, ok := omittedValues[value]; ok {
					omitted++
					continue
				}

				deck = append(deck, Card{suit, value})
			}
		}

		for j := 0; j < conf.jokers; j++ {
			var suit Suit

			if j&1 == 1 {
				suit = Red
			} else {
				suit = Black
			}

			deck = append(deck, Card{suit, Joker})
		}

		decks = append(decks, deck...)
	}

	return conf.sort(decks)
}
