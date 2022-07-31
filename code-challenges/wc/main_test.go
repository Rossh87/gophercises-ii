package main

import "testing"

func TestCount(t *testing.T) {
	gotSingleChar := count("s")

	wantFromSingleChar := 1

	if gotSingleChar != wantFromSingleChar {
		t.Fatalf("single char failed")
	}

	gotSingleWord := count("single")

	wantFromSingleWord := 1

	if gotSingleWord != wantFromSingleWord {
		t.Fatalf("single word failed")
	}

	gotMultiWord := count("moreWords")

	wantFromMultiWord := 2

	if gotMultiWord != wantFromMultiWord {
		t.Fatalf("multiword failed")
	}
}
