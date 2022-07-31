package main

import (
	"testing"
)

func TestSolicitAnswer(t *testing.T) {
	var answer = make(chan string)

	SolicitAnswer(answer)

}
