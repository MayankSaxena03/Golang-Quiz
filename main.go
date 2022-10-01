package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

func problemPuller(fileName string) ([]problem, error) {
	//Open the file
	fObj, err := os.Open(fileName)
	if err == nil {
		//We will create a new reader
		csvR := csv.NewReader(fObj)
		//Read the lines
		cLines, err := csvR.ReadAll()
		if err == nil {
			//Call Parse Problem
			return parseProblem(cLines), nil
		}
	}
	return nil, err
}

func main() {
	//Input Name of the File
	var fName string
	fName = "/quiz.csv"
	//Set Duration of the timer
	timer := 30
	//Pull the questions from the file
	problems, err := problemPuller(fName)
	//Handle The Error
	if err != nil {
		panic(err)
	}
	//Count correct answers
	correctAns := 0
	//Initialize the timer using duration of the timer
	tObj := time.NewTimer(time.Duration(timer) * time.Second)
	ansC := make(chan string)
	//Loop through the problems and input the answers
problemLoop:
	for i, p := range problems {
		var ans string
		fmt.Printf("\nProblem #%d: %s = ", i+1, p.q)
		go func() {
			fmt.Scanf("%s ", &ans)
			ansC <- ans
		}()
		select {
		case <-tObj.C:
			fmt.Println("Time is up!")
			break problemLoop
		case IAns := <-ansC:
			if IAns == p.a {
				// fmt.Printf("Correct Answer! %s\n", p.a)
				correctAns++
			}
			if i == len(problems)-1 {
				close(ansC)
			}
		}
	}
	//Print the Results
	fmt.Printf("\nYou Scored %d out of %d\n", correctAns, len(problems))
	fmt.Println("Press any key to exit...")
	<-ansC
}

func parseProblem(lines [][]string) []problem {
	//Go over the lines and create a problem with problem struct
	r := make([]problem, len(lines))
	for i := 0; i < len(lines); i++ {
		r[i] = problem{q: lines[i][0], a: lines[i][1]}
	}
	return r
}

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
