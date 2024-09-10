package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	timerCh := make(chan string)
	var (
		userAnswer                                        string
		totalQuestions, answeredQuestions, correctAnswers int
	)

	// defining flags
	csvFilePath := flag.String("f", "problems.csv", "set csv('question,answer') file path; default: 'problems.csv'")
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
	fmt.Println("This is a Quiz program written in GOas a task from 'https://github.com/gophercises/quiz'.\n-----")

	// TODO: making order of quiz questions

	// Wait user to start quiz
	fmt.Println("Press ENTER when you're ready: ")
	reader := bufio.NewReader(os.Stdin)
	_, errR := reader.ReadString('\n')
	if errR != nil {
		log.Fatal(errR)
	}

	// Start timer
	go func() {
		// first convert timer seconds given in -s flag(30 default) to string to parse time duration
		timerSecondsString := strconv.Itoa(*timerSeconds)

		// setting duration for timer
		duration, err := time.ParseDuration(timerSecondsString + "s")
		if err != nil {
			log.Fatalf("FAILED: to parse duration from user input:\n\t%v", err)
		}

		// starting timer, sending msg to channel when expired
		timer := time.NewTimer(duration)
		<-timer.C
		timerCh <- "Timer is EXPIRED!"
	}()

	// quiz loop
	fmt.Println("Quiz is started!")
	for _, v := range problems {
		select {
		case msg := <-timerCh:
			fmt.Println(msg)
			fmt.Printf("\nYour result is %d correct answers of %d answered of %d total questions!\n", correctAnswers, answeredQuestions, totalQuestions)
			return
		default:
			fmt.Printf("The question is: %s\n", v.question)
			fmt.Print("Your answer: ")
			fmt.Scan(&userAnswer)

			if *&userAnswer == v.answer {
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
