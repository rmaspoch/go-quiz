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

func init() {
	flag.IntVar(&timeLimit, "limit", 15, "Time limit in seconds")
}

func main() {
	// read time limit from command line
	flag.Parse()
	fmt.Printf("Time limit set to %v seconds\n", timeLimit)

	quiz := readQuiz("quiz.csv")
	totalQuestions := len(quiz)
	correctAnswers := 0

	correct := make(chan struct{})
	done := make(chan bool)

	go func() {
		runQuiz(quiz, correct, done)
	}()

	exitLoop := false
	for !exitLoop {
		select {
		case <-done:
			fmt.Println("Done")
			exitLoop = true
		case <-correct:
			correctAnswers++
		case <-time.After(time.Duration(timeLimit) * time.Second):
			fmt.Println("\nTime limit reached")
			exitLoop = true
		}
	}

	fmt.Printf("You answered %v questions correctly out of %v", correctAnswers, totalQuestions)
}

func readQuiz(fileName string) (quiz []questionAnswer) {
	file, err := os.Open(fileName)
	guard(err)

	reader := csv.NewReader(file)
	var lines [][]string
	lines, err = reader.ReadAll()
	guard(err)

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
func runQuiz(quiz []questionAnswer, correct chan<- struct{}, done chan<- bool) {
	for _, q := range quiz {
		fmt.Printf("How much is %v ", q.question)
		var answer string
		fmt.Scanf("%s\n", &answer)
		if answer == q.answer {
			correct <- struct{}{}
		} else {
			fmt.Println("Incorrect")
		}
	}
	done <- true
}

func guard(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
