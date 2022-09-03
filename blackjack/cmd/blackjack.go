package main

import (
	"flag"
	"fmt"
	"gophercises/blackjack/blackjack"
	"os"
)

type logger struct {
}

func (l logger) Print(s string) {
	fmt.Print(s)
}

func main() {
	var players int

	flag.IntVar(&players, "p", 1, "set number of players")

	flag.Parse()

	logger := logger{}

	blackjack.Run(logger, players, os.Stdin)
}
