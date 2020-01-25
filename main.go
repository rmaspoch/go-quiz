package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type questionAnswer struct {
	question string
	answer   string
}

// time limit flag
var timeLimit int

func main() {
	flag.IntVar(&timeLimit, "limit", 15, "Time limit in seconds")
	flag.Parse()
	fmt.Printf("Time limit set to %v seconds\n", timeLimit)

	quiz := readQuiz("quiz.csv")
	totalQuestions := len(quiz)

	correctAnswers := runQuiz(quiz)
	fmt.Printf("You answered %v questions correctly out of %v\n", correctAnswers, totalQuestions)
}

func readQuiz(fileName string) (quiz []questionAnswer) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	reader := csv.NewReader(file)
	var lines [][]string
	lines, err = reader.ReadAll()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return parseLines(lines)
}

func parseLines(lines [][]string) []questionAnswer {
	quiz := make([]questionAnswer, len(lines))
	for i, record := range lines {
		quiz[i] = questionAnswer{
			question: record[0],
			answer:   strings.TrimSpace(record[1]),
		}
	}
	return quiz
}

// returns number of correct answers
func runQuiz(quiz []questionAnswer) int {
	correctAnswers := 0
	timeout := time.NewTimer(time.Duration(timeLimit) * time.Second)
	result := make(chan string)

	for i, q := range quiz {
		fmt.Printf("Question #%d - %v ", i+1, q.question)

		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			result <- answer
		}()

		select {
		case <-timeout.C:
			fmt.Println("\nYou have run out of time!")
			return correctAnswers
		case answer := <-result:
			if answer == q.answer {
				correctAnswers++
			} else {
				fmt.Println("Incorrect!")
			}
		}
	}
	return correctAnswers
}
