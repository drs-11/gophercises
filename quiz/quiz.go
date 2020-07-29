package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func getRecords(f string) [][]string {
	file, err := os.Open(f)
	if err != nil {
		fmt.Println("Something bad happened while opening file.")
		log.Fatal(err)
	}
	defer file.Close()

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	return records

}
func main() {
	quizFile := flag.String("csv", "problems.csv", "a csv file which contains questions and answers")
	timeLimit := flag.Int("limit", 30, "time limit for quiz")
	flag.Parse()

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	records := getRecords(*quizFile)
	total := len(records)
	var correct, wrong, i int

loop:
	for {
		ansCh := make(chan string)

		go func() {
			var ans string
			fmt.Printf("Q: %s = ", records[i][0])
			fmt.Scanln(&ans)
			ansCh <- ans
		}()

		select {
		case <-timer.C:
			fmt.Println("\nTime up!")
			break loop
		case answer := <-ansCh:
			if i >= total {
				break loop
			}

			if answer == records[i][1] {
				correct++
			} else {
				wrong++
			}
			i++
		}

	}

	fmt.Println("You got", correct, "correct and", wrong, "wrong out of", total)

}
