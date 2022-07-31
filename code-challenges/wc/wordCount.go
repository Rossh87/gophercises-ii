package main

import (
	"fmt"
	"strings"
)

func main() {
	var input string

	fmt.Scanf("%s\n", &input)

	ret := count(input)

	fmt.Printf("%d\n", ret)
}

func isUpper(r rune) (result bool) {
	s := string(r)
	result = strings.ToUpper(string(s)) == string(s)
	return
}

func count(s string) int {
	count := 1

	for _, rune := range s {
		if isUpper(rune) {
			count++
		}
	}

	return count
}
