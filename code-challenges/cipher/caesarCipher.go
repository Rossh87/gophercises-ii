package main

import (
	"fmt"
	"strings"
)

func main() {
	var input string
	var rotateCount int

	fmt.Println("input:")
	fmt.Scanf("%s\n", &input)

	fmt.Println("rotation count:")
	fmt.Scanf("%d\n", &rotateCount)

	ret := cipher(input, rotateCount)

	fmt.Printf("%s\n", ret)
}

func getCharAtPos(s string, pos int) rune {
	var empty rune

	ct := -1

	for _, r := range s {
		ct++
		if ct == pos {
			return rune(r)
		}
	}

	return empty
}

func rotate(r rune, rotation int, alphabet string) rune {
	initPos := strings.IndexRune(alphabet, r)

	if initPos < 0 {
		panic("rune not in alphabet!")
	}

	pos := (initPos + rotation) % len(alphabet)

	return getCharAtPos(alphabet, pos)
}

func cipher(s string, rotation int) string {
	alphabet := `abcdefghijklmnopqrstuvwxyz`
	upperAlphabet := `ABCDEFGHIJKLMNOPQRSTUVWXYZ`

	var ret string

	for _, r := range s {
		switch {
		case strings.ContainsRune(alphabet, r):
			ret += string(rotate(r, rotation, alphabet))

		case strings.ContainsRune(upperAlphabet, r):
			ret += string(rotate(r, rotation, upperAlphabet))

		default:
			ret += string(r)
		}
	}

	return ret
}
