package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type problem struct {
	question string
	answer   int
}

func main() {
	var csvPath string
	flag.StringVar(&csvPath, "csv", "problems.csv", "Provide the csv filepath")
	var timeoutSeconds int
	flag.IntVar(&timeoutSeconds, "timeout", 5, "Provide the question timeout in seconds")
	flag.Parse()

	file, err := os.Open(csvPath)
	if err != nil {
		log.Fatal("Error while reading the file", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	if err != nil {
		fmt.Println("Error reading records")
	}

	fmt.Printf("You have %d seconds to answer each question, are you ready? (Hit Enter)", timeoutSeconds)
	fmt.Scanln()

	score := 0
	done := make(chan bool)

	go func() {
		for _, record := range records {
			truth, _ := strconv.Atoi(record[1])
			curr := problem{question: record[0], answer: truth}

			var userAnswer int
			fmt.Println(curr.question, "?")
			fmt.Scanf("%d", &userAnswer)

			if userAnswer == curr.answer {
				score++
			}
		}
		done <- true
	}()

	select {
	case <-done:
		fmt.Println("You finished within time!")
	case <-time.After(time.Duration(timeoutSeconds) * time.Second):
		fmt.Println("Timeout reached!")
	}

	fmt.Println("Your total score:", score, "/", len(records))
}
