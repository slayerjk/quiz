package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"slices"
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
	isShuffled := flag.Bool("s", false, "shuffle quiz file questions order; default: false")
	flag.Parse()

	// parsing csv file
	quizData, err := parseCsv(*&csvFilePath)
	if err != nil {
		log.Fatalf("failed to parse csv file %s:\n\t%v", *csvFilePath, err)
	}
	totalQuestions = len(quizData)

	// printing app intro info & options selected
	fmt.Println("This is a Quiz program written in GO as a task from 'https://github.com/gophercises/quiz'.\n-----")

	// making order of quiz questions
	quizOrder := defineQuizOrder(*isShuffled, len(quizData))
	fmt.Println(quizOrder)

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
	for _, v := range quizOrder {
		fmt.Printf("The question is: %s\n", quizData[v].question)

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
			if strings.TrimSpace(strings.ToLower(answer)) == quizData[v].answer {
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
			// normalizing correct answer
			answer: strings.TrimSpace(strings.ToLower(v[1])),
		}
	}

	return result, nil
}

func defineQuizOrder(flag bool, lenOfData int) []int {
	result := make([]int, 0)
	randNum := rand.Intn(lenOfData)

	if !flag {
		for i := 0; i < lenOfData; i++ {
			result = append(result, i)
		}
	} else {
		for i := 0; i < lenOfData; i++ {
			for {
				if isPresent := slices.Contains(result, randNum); isPresent {
					randNum = rand.Intn(lenOfData)
					continue
				} else {
					break
				}
			}
			result = append(result, randNum)
		}
	}

	return result
}
