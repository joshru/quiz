package main

import (
	"flag"
	"os"
	"fmt"
	"encoding/csv"
	"strings"
	"time"
)

// flags
var (
	csvFilename = flag.String( "csv", "problems.csv", "a csv file in the format of 'question,answer'" )
	timeLimit   = flag.Int("limit", 3, "time limit for the quiz in seconds")
)

func main() {

	flag.Parse()

	file, err := os.Open( *csvFilename )
	if err != nil {
		exit(fmt.Sprintf("Failed ot open the CSV file: '%s'\n", *csvFilename))
	}
	defer file.Close()

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the provided CSV file.")
	}

	problems := parseLines(lines)
	timer := time.NewTimer(time.Second * time.Duration(*timeLimit))
	correct := 0
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.q)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <- timer.C:
			fmt.Printf("\nYou scored %d out of %d.\n", correct, len(problems))
			return

		case answer := <-answerCh:
			if answer == p.a {
				correct++
			}
		}
	}
}

type problem struct {
	q string
	a string
}

func parseLines(lines [][]string) []problem {
	ret := make([] problem, len(lines))

	for i, line := range lines {
		ret[i] = problem {
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

func exit (msg string) {
	fmt.Println(msg)
	os.Exit(1)
}