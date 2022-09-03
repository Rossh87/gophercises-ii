package blackjack

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

type Game struct {
	players []Player
	dealer  []Card
	deck    Deck
}

func (g Game) Player(id int) (*Player, error) {
	for _, p := range g.players {
		if p.id == id {
			return &p, nil
		}
	}

	return nil, errors.New("requested a player id that is not registered with the game")
}

func (g Game) Total() int {
	return total(g.dealer)
}

func (g Game) Done() bool {
	doneCount := 0

	for _, p := range g.players {
		if p.stands || p.bust {
			doneCount++
		}
	}

	return doneCount == len(g.players)
}

func NewGame(playerCount int) *Game {
	deck := NewDeck(WithDecks(3))

	players := make([]Player, playerCount)

	for i := range players {
		h := []Card{}
		players[i] = Player{
			id:     i + 1,
			hand:   h,
			stands: false,
			bust:   false,
		}
	}

	return &Game{
		players,
		[]Card{},
		deck,
	}
}

func (g *Game) Deal() {
	for j := 0; j < 2; j++ {
		for i := 0; i <= len(g.players); i++ {
			c, _ := (g.deck.Deal())

			if i == len(g.players) {
				g.dealer = append(g.dealer, *c)
			} else {
				g.players[i].hand = append(g.players[i].hand, *c)
			}
		}
	}
}

type Logger interface {
	Print(s string)
}

func Run(logger Logger, playerCount int, input *os.File) {
	scanner := bufio.NewScanner(input)

	logger.Print("Beginning a new game!\n")

	g := NewGame(playerCount)

	g.Deal()

	for _, p := range g.players {
		logger.Print(fmt.Sprintf("Player %d is dealt\n%s", p.id, p.PrintHand()))
	}

	logger.Print(fmt.Sprintf("Dealer is dealt\n%s", printhand(g.dealer, &[]int{0})))

	p := -1

	for {
		if g.Done() {
			break
		}

		p++

		player := &g.players[p%len(g.players)]

		if player.bust || player.stands {
			continue
		}

		logger.Print(fmt.Sprintf("Player %d, your hand is\n%s", player.id, player.PrintHand()))
		logger.Print("Do you want to hit (y/n)?")

		var hit string

		scanner.Scan()

		hit = scanner.Text()

		if hit == "y" {
			c, _ := g.deck.Deal()
			logger.Print(fmt.Sprintf("Player %d receives a %s", player.id, c.Print()))
			player.hand = append(player.hand, *c)
		} else if hit == "n" {
			logger.Print(fmt.Sprintf("Player %d stands\n", player.id))
			player.stands = true
		} else {
			panic("received invalid input in response to hit question!!")
		}

		if player.Total() > 21 {
			logger.Print(fmt.Sprintf("Player %d is bust!\n", player.id))
			player.bust = true
		}
	}

	bustCount := 0

	for _, p := range g.players {
		if p.bust {
			bustCount++
		}
	}

	if bustCount == len(g.players) {
		logger.Print("All players busted, dealer wins! Better luck next time assholes.")
		os.Exit(0)
	} else {
		logger.Print(fmt.Sprintf("Dealer's hand:\n%s", printhand(g.dealer, nil)))

		dealerTotal := total(g.dealer)

		maxPlayerTotal := 0

		currMaxPlayer := 0

		for _, p := range g.players {
			t := p.Total()

			if t > maxPlayerTotal {
				maxPlayerTotal = t
				currMaxPlayer = p.id
			}
		}

		if maxPlayerTotal > dealerTotal {
			winner, _ := g.Player(currMaxPlayer)

			logger.Print(fmt.Sprintf("Congratulations player %d, you win with a total of %d!\n%s", currMaxPlayer, maxPlayerTotal, winner.PrintHand()))
			os.Exit(0)
		} else {
			logger.Print(fmt.Sprintf("Dealer wins with a total of %d", dealerTotal))
			logger.Print(fmt.Sprintf("Dealer's hand:\n%s", printhand(g.dealer, nil)))
			logger.Print("Better luck next time!")
			os.Exit(0)
		}
	}
}
