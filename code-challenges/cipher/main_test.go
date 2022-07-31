package main

import "testing"

type testCase struct {
	rotation int
	input    string
	want     string
}

func TestCipher(t *testing.T) {
	cases := []testCase{
		{
			0, "a", "a",
		},
		{
			1, "z", "a",
		},
		{
			1, "A-z", "B-a",
		},
		{
			2, "A*k-B", "C*m-D",
		},
	}

	for _, tc := range cases {

		got := cipher(tc.input, tc.rotation)

		if got != tc.want {
			t.Fatalf("Wanted: %s\nGot: %s", tc.want, got)
		}
	}
}
