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
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of question, answer ")
	timeLimit := flag.Int("limit", 30, "timelimit of the quiz in seconds")
	flag.Parse()
	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("failed to open CSV file: %s", *csvFilename))
	}
	_ = file

	r := csv.NewReader(file)
	lines, err := r.ReadAll()

	if err != nil {
		exit("Failed to parse provided CSV file")
	}

	// fmt.Println(lines)
	problems := parseLines(lines)
	fmt.Println(problems)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	// <-timer.C

	totalscore := 0

	for i, problem := range problems {
		fmt.Printf("Problem # %d: %s = \n", i+1, problem.q)

		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("You scored %d correct out of %d \n", totalscore, len(problems))
			return
		case answer := <-answerCh:
			if answer == problem.a {
				totalscore++
			}
		}
	}

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

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
