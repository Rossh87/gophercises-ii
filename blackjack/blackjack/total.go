package blackjack

import (
	"strings"
)

func total(hand []Card) int {
	aces := 0

	total := 0

	for _, c := range hand {
		if c.Value == Ace {
			aces++
		}

		total += c.Value.Value()
	}

	for i := 0; i < aces; i++ {
		// try to use 11.  If adding 11 puts us over total, add 1 instead.
		// Account for the possibility there is more than 1 ace in players hand.
		target := 21 - maxInt(0, aces-i-1)

		if total+11 <= target {
			total += 11
		} else {
			total += 1
		}
	}

	return total
}

func contains(a int, as []int) bool {
	for _, el := range as {
		if el == a {
			return true
		}
	}

	return false
}

func printhand(hand []Card, hiddenCards *[]int) string {
	var s strings.Builder

	for i, c := range hand {
		var subStr string

		if hiddenCards != nil && contains(i, *hiddenCards) {
			subStr = "*\n"
		} else {
			subStr = c.Print()
		}

		s.WriteString(subStr)
	}

	return s.String()
}
