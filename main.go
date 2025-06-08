package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

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

func startQuiz() {
	fmt.Println("Press ENTER to start the quiz, any other key will terminate the program")

	_, key, err := keyboard.GetSingleKey()
	if err != nil {
		log.Fatal(err)
	}

	if key != keyboard.KeyEnter {
		log.Fatal("Terminating program")
	}
}

func main() {
	var correct int

	fname := flag.String("fname", "problems.csv", "specify the filename to open")
	flag.Parse()

	f, err := os.Open(*fname)
	if err != nil {
		log.Fatal(err)
	}

	questions := parseCsv(f)
	startQuiz()

	for _, q := range questions {
		var answer string

		fmt.Printf("%s = ", q[0])
		fmt.Scanln(&answer)

		if answer == q[1] {
			correct++
		}
	}

	fmt.Printf("\nTotal questions: %d\nCorrect answers: %d\n", len(questions), correct)
}
