package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
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
		var answer string

		results <- fmt.Sprintf("%s = ", q[0])
		fmt.Scanln(&answer)

		if answer == q[1] {
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
