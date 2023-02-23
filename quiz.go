package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type QA struct {
	question string
	answer   string
}

func main() {
	file, timeout, err := getParams()
	if err != nil {
		log.Fatal(err)
		return
	}
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer f.Close()
	content := parseCsv(readCsv(f))
	score := 0
	fmt.Printf("Press enter to start the timer (%d)...\n", timeout)
	for {
		if takeUserAnswer() != "" {
			break
		}
	}
	timer1 := time.NewTimer(time.Duration(timeout) * time.Second)
	go func() {
		<-timer1.C
		fmt.Printf("Time is up!")
		showScore(score)
		os.Exit(0) // probably bullshit method to exit
	}()

	// main loop for verifying the answers
	for i, val := range content {
		fmt.Printf("Question %d: %s\n", i+1, val.question)
		answer := strings.TrimSpace(takeUserAnswer())
		if answer == strings.TrimSpace(val.answer) {
			fmt.Println("Correct!")
			score += 1
		} else {
			fmt.Println("Wrong!")
		}
	}
	showScore(score)

}

func takeUserAnswer() string {
	reader := bufio.NewReader(os.Stdin)
	answer, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}
	return answer
}

func getParams() (string, int, error) {
	filePtr := flag.String("file", "problems.csv", "A file used for the quiz")
	timePtr := flag.Int("timer", 20, "A time in seconds to finish the quiz")
	flag.Parse()
	fmt.Printf("File passed is %s\n", *filePtr)
	if strings.HasSuffix(*filePtr, ".csv") == false {
		return "", *timePtr, errors.New("Not a csv file!\n")
	}
	return *filePtr, *timePtr, nil
}

func readCsv(f *os.File) [][]string {
	reader := csv.NewReader(f)
	all, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return all
}

func parseCsv(cont [][]string) []QA {
	allStuff := make([]QA, len(cont))
	for i, val := range cont {
		var qa QA
		qa.question = val[0]
		qa.answer = val[1]
		allStuff[i] = qa
	}
	return allStuff
}
func showScore(score int) {
	fmt.Printf("Bravo, final score is %d", score)
}
