package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
)

type question []string

func parseCsv(r io.Reader) (q []question) {
	reader := csv.NewReader(r)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		q = append(q, record)
	}

	return q
}

func startQuiz(questions []question, results chan string, correct *int) {
	for _, q := range questions {
		var input string
		question := q[0]
		answer := q[1]

		results <- fmt.Sprintf("%s = ", question)
		fmt.Scanln(&input)

		answer = strings.Trim(answer, " ")
		answer = strings.ToLower(answer)
		input = strings.ToLower(input)

		if input == answer {
			*correct++
		}
	}

	close(results)
}

func initializeTimeout(duration *int) (timeout <-chan time.Time) {
	fmt.Println("Press ENTER to start the quiz, any other key will terminate the program")

	_, key, err := keyboard.GetSingleKey()
	if err != nil {
		log.Fatal(err)
	}

	if key != keyboard.KeyEnter {
		log.Fatal("Terminating program")
	}

	return time.After(time.Second * time.Duration(*duration))
}

func main() {
	var (
		correct int
		runtime = flag.Int("time", 30, "specify the duration of the quiz in seconds")
		fname   = flag.String("fname", "problems.csv", "specify the filename to open")
		// shuffle = flag.Bool("rand", true, "whether or not shuffle quiz questions")
	)

	flag.Parse()

	f, err := os.Open(*fname)
	if err != nil {
		log.Fatal(err)
	}

	questions := parseCsv(f)
	results := make(chan string, len(questions))

	timeout := initializeTimeout(runtime)
	go startQuiz(questions, results, &correct)

loop:
	for {
		select {
		case <-timeout:
			fmt.Println("\nTime is up")
			break loop
		case v, ok := <-results:
			if !ok {
				break loop
			}

			fmt.Printf("%v", v)
		}
	}

	fmt.Printf("\nCorrect answers: %d\nTotal questions: %d\n", correct, len(questions))
}
