package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	csvFilePath *string
	limit       *int
)

func init() {
	csvFilePath = flag.String("csv", "./questions.csv", "Path to csv file with problems")
	limit = flag.Int("limit", 30, "The quiz duration limit in seconds")
}

func main() {
	flag.Parse()
	fmt.Printf("The limit is set to %v seconds\n", *limit)
	file, err := os.Open(*csvFilePath)
	if err != nil {
		exit(fmt.Sprintf("Failed to open file %v", *csvFilePath))
	}

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		exit("Something went wrong with parsing the problems")
	}

	timer := time.NewTimer(time.Duration(*limit) * time.Second)

	answerChan := make(chan string)

	problems := parseLines(lines)
	score := 0

loop:
	for _, problem := range problems {
		fmt.Printf("%v = ", problem.q)
		// go routine, make scanning non blocking
		go scanForAnswer(answerChan)

		select {
		// receive "notification" from timer
		case <-timer.C:
			break loop
		// receive an answer from user through channel
		case ans := <-answerChan:
			if ans == problem.a {
				score++
			}
		}
	}

	fmt.Printf("\nYou scored %v out of %v\n", score, len(problems))
}

type problem struct {
	q string
	a string
}

func parseLines(lines [][]string) []problem {
	var result = make([]problem, len(lines))
	for i, line := range lines {
		result[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}

	return result
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func scanForAnswer(c chan string) {
	var userAnswer string
	fmt.Scanf("%s", &userAnswer)
	c <- userAnswer
}
