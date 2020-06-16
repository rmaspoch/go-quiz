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

const (
	defaultLimit    = 15
	defaultQuizFile = "quiz.csv"
)

func main() {
	// time limit flag
	var timeLimit int
	// file flag
	var quizFile string

	flag.IntVar(&timeLimit, "limit", defaultLimit, "Time limit in seconds")
	flag.StringVar(&quizFile, "file", defaultQuizFile, "CSV file with questions and answers")
	flag.Parse()

	fmt.Printf("Using %s CSV file\n", quizFile)
	fmt.Printf("Time limit set to %v seconds\n", timeLimit)

	quiz := readQuiz(quizFile)
	totalQuestions := len(quiz)

	correctAnswers := runQuiz(quiz, timeLimit)
	fmt.Printf("You answered %v questions correctly out of %v\n", correctAnswers, totalQuestions)
}

func readQuiz(fileName string) (quiz []questionAnswer) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Unable to open file: %v\n", err)
	}

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Unable to read file: %v\n", err)
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
func runQuiz(quiz []questionAnswer, timeLimit int) int {
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
