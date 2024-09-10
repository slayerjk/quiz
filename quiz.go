package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	var (
		userAnswer                                        string
		totalQuestions, answeredQuestions, correctAnswers int
	)

	// defining flags
	csvFilePath := flag.String("f", "problems.csv", "set csv('question,answer' format) file path; default: 'problems.csv'")
	timerSeconds := flag.Int("t", 30, "time execution limit, seconds; default: 30")
	// isShuffled := flag.Bool("s", false, "shuffle quiz file questions order; default: false")
	flag.Parse()

	// parsing csv file
	problems, err := parseCsv(*&csvFilePath)
	if err != nil {
		log.Fatalf("failed to parse csv file %s:\n\t%v", *csvFilePath, err)
	}
	totalQuestions = len(problems)

	// printing app intro info & options selected
	fmt.Println("This is a Quiz program written in GO as a task from 'https://github.com/gophercises/quiz'.\n-----")

	// TODO: making order of quiz questions

	// Wait user to start quiz
	fmt.Println("Press ENTER when you're ready: ")
	reader := bufio.NewReader(os.Stdin)
	_, errR := reader.ReadString('\n')
	if errR != nil {
		log.Fatal(errR)
	}

	// starting timer, sending msg to channel when expired
	timer := time.NewTimer(time.Duration(*timerSeconds) * time.Second)

	// quiz loop
	fmt.Println("Quiz is started!")

quizLoop:
	for _, v := range problems {
		fmt.Printf("The question is: %s\n", v.question)

		// channel to recieve user answers
		answerCh := make(chan string)

		// running goroutine to get & validate user answers
		go func() {
			fmt.Print("Your answer: ")
			fmt.Scan(&userAnswer)

			answerCh <- userAnswer
		}()

		// waiting either for time expiration or continue to ask questions
		select {
		case <-timer.C:
			fmt.Println("Timer is expired!")
			break quizLoop
		case answer := <-answerCh:
			// compare normalized user answer with correct answer
			if strings.TrimSpace(strings.ToLower(answer)) == v.answer {
				correctAnswers++
			}
			answeredQuestions++
		}
	}

	fmt.Printf("\nYour result is %d correct answers of %d answered of %d total questions!\n", correctAnswers, answeredQuestions, totalQuestions)
}

type quiz struct {
	question string
	answer   string
}

func parseCsv(path *string) ([]quiz, error) {
	file, err := os.Open(*path)
	if err != nil {
		return nil, fmt.Errorf("failed to open csv file %s:\n\t%v", *path, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	csvData, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read csv file %s:\n\t%v", *path, err)
	}

	result := make([]quiz, len(csvData))

	for i, v := range csvData {
		result[i] = quiz{
			question: v[0],
			answer:   strings.TrimSpace(v[1]),
		}
	}

	return result, nil
}
