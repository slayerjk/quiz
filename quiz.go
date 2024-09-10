package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type Quiz struct {
	question string
	answer   string
}

func main() {
	timerCh := make(chan string)
	quiz := Quiz{}
	var (
		quizeList                                         []Quiz
		userAnswer, cleanedAnswer                         string
		totalQuestions, answeredQuestions, correctAnswers int
		csvFilePath                                       string = "problems.csv"
		timerSeconds                                      string = "30s"
		optionsSummary                                    string
		appArgs                                           = make(map[string]string)
	)
	helpMsg := `This is a Quiz program written in GO as a task from "https://github.com/gophercises/quiz".

Help:

NO ARGS: run program with default parameters:
	- default "problem.csv" quiz questions csv file located in the same directorty as the quiz main program;
	- 30 secs for execution (countdown timer);
	- no shuffle for quiz questions 

-h: print the help info and EXIT program
-f FILEPATH: use different .csv file(must have two columns without header: "question,answer")
-t SECONDS: custom countdown timer, in seconds; can be used with "-f" separately
-s: turn Shuffle of quiz questions ON

Examples:
	quiz # run program with default ./problems.csv and 30 seconds for timer
	quiz -f test/test.csv # run program with test/test.csv as csv input file
	quiz -f test/test.csv -t 60 # run program with test/test.csv as csv input file and 60 secs for timer
	quiz -t 60 # run program with 60 secs for timer
	quiz -t 60 -s # run program with 60 secs for timer and turn shuffle on
-----`

	// getting all args
	for i, v := range os.Args {
		if i == 0 {
			continue
		} else {
			switch {
			case v == "-h":
				fmt.Println(helpMsg)
				os.Exit(0)
			case v == "-s":
				appArgs[v] = ""
				optionsSummary += "Shuffle quiz questions is: On\n"
			case v == "-f" || v == "-t":
				if os.Args[i+1] != "" {
					if v == "-f" {
						appArgs[v] = ""
						csvFilePath = os.Args[i+1]
						optionsSummary += fmt.Sprintf("Using %s as csv file path\n", csvFilePath)
					} else if v == "-t" {
						appArgs[v] = "ON"
						timerSeconds = os.Args[i+1] + "s"
						optionsSummary += fmt.Sprintf("Using %s for timer\n", timerSeconds)
					} else {
						log.Fatalf("Empty arg for %s", v)
					}
				}
			default:
				continue
			}
		}
	}

	// check if shuffle is off for default message
	if _, ok := appArgs["-s"]; !ok {
		optionsSummary += "Shuffle quiz questions is: Off(Default)\n"
	}
	// check if -f options true
	if _, ok := appArgs["-f"]; !ok {
		optionsSummary += fmt.Sprintf("Using Default %s as csv file path\n", csvFilePath)
	}
	// check if -t options true
	if _, ok := appArgs["-t"]; !ok {
		optionsSummary += fmt.Sprintf("Using %s for timer\n", timerSeconds)
	}

	// reading quiz csv file
	csvFile, err := os.Open(csvFilePath)
	if err != nil {
		log.Fatalf("failed to open quiz csv file:\n\t%v", err)
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	csvReader.FieldsPerRecord = 2

	csvData, err := csvReader.ReadAll()
	if err != nil {
		log.Fatalf("failed to read from csvReader:\n\t%v", err)
	}

	// printing app intro info & options selected
	fmt.Println(`This is a Quiz program written in GOas a task from "https://github.com/gophercises/quiz".

	Use "-h" flag to print help.
-----`)
	fmt.Println(optionsSummary)

	// parsing and checking csv file
	for _, v := range csvData {
		quiz.question = v[0]
		quiz.answer = v[1]

		// answer or question in csv can't be empty
		if quiz.question == "" || quiz.answer == "" {
			log.Fatalf("quiz answer or question can't be empty,\n\t check this entry: %v", v)
		}
		quizeList = append(quizeList, quiz)
		totalQuestions++
	}

	// making order of quiz questions
	indexes := make([]int, 0)
	println(indexes)

	// if -s is ON then make slice of random unique indexes
	if appArgs["-s"] != "ON" {
		for i := 0; i < len(quizeList); i++ {
			indexes = append(indexes, i)
		}
		fmt.Println(indexes)
	} else {

	}
	os.Exit(0)

	// Wait user to start quiz
	fmt.Println("Press ENTER when you're ready: ")
	reader := bufio.NewReader(os.Stdin)
	_, errR := reader.ReadString('\n')
	if errR != nil {
		log.Fatal(errR)
	}

	// Start timer
	go func() {
		duration, err := time.ParseDuration(timerSeconds)
		if err != nil {
			log.Fatalf("FAILED: to parse duration from user input:\n\t%v", err)
		}

		timer := time.NewTimer(duration)
		<-timer.C
		timerCh <- "Timer is EXPIRED!"
	}()

	// quiz loop
	fmt.Println("Quiz is started!")
	for _, v := range quizeList {
		select {
		case msg := <-timerCh:
			fmt.Println(msg)
			fmt.Printf("\nYour result is %d correct answers of %d answered of %d total questions!\n", correctAnswers, answeredQuestions, totalQuestions)
			return
		default:
			fmt.Printf("The question is: %s\n", v.question)
			fmt.Print("Your answer: ")
			fmt.Scan(&userAnswer)

			// clean up an answer(trim trailing whitespaces, convert to lower)
			cleanedAnswer = strings.ToLower(strings.Trim(userAnswer, " "))
			if cleanedAnswer == v.answer {
				correctAnswers++
			}
			answeredQuestions++
		}
	}
	fmt.Printf("\nYour result is %d correct answers of %d answered of %d total questions!\n", correctAnswers, answeredQuestions, totalQuestions)
}
