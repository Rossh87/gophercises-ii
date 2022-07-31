package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")

	timeLimit := flag.Int("timer", 30, "Specify the time limit, in seconds, the user will have to complete quiz")

	flag.Parse()

	file, err := os.Open(*csvFilename)

	if err != nil {
		exitStr := fmt.Sprintf("failed to open the CSV file %s\n", *csvFilename)
		exit(exitStr)
	}

	r := csv.NewReader(file)

	lines, err := r.ReadAll()

	if err != nil {
		exitStr := fmt.Sprintf("failed to read the CSV file %s\n", *csvFilename)
		exit(exitStr)
	}

	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0

	// try breaking this into its own function and writing a test!
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, p.q)

		answerCh := make(chan string)

		SolicitAnswer(answerCh)

		select {
		case <-timer.C:
			fmt.Printf("You scored %d out of %d\n", correct, len(problems))
			return
		case answer := <-answerCh:
			if answer == p.a {
				correct++
				fmt.Println("Correct!")
			}
		}
	}
	fmt.Printf("You scored %d out of %d\n", correct, len(problems))

}

func SolicitAnswer(receiver chan string) {
	go func() {
		var ans string
		fmt.Scanf("%s\n", &ans)
		receiver <- ans
	}()
}
func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))

	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}

	return ret
}

func exit(msg string) {
	fmt.Println(msg)

	os.Exit(1)
}

// better to depend on this type rather than the 2D slice we get back from parsing CSV.  This way,
// in the future our program could parse different kinds of input, and the rest of the code would work fine
// provided we manipulate the data into this shape first.
type problem struct {
	q string
	a string
}
