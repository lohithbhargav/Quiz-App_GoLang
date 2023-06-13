package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

type problem struct {
	question string
	answer   string
}

func loadQuestionFile(fileName string) ([]problem, error) {
	if fileObj, err := os.Open(fileName); err == nil {
		csvReader := csv.NewReader(fileObj)
		if lines, err := csvReader.ReadAll(); err == nil {
			return parseProblem(lines), nil
		} else {
			return nil, fmt.Errorf("error in reading data in csv %s %s", fileName, err.Error())
		}
	} else {
		return nil, fmt.Errorf("error in opening csv %s %s", fileName, err.Error())
	}
}

func parseProblem(lines [][]string) []problem {
	r := make([]problem, len(lines))
	for i := 0; i < len(lines); i++ {
		r[i] = problem{question: lines[i][0], answer: lines[i][1]}
	}
	return r
}

func main() {
	//fileName := flag.String("f", "quiz.csv", "path of csv file")
	fileName := "quiz.csv"
	//timer := flag.Int("t", 30, "timer for the quiz")
	timer := 30
	//flag.Parse()
	problems, err := loadQuestionFile(fileName)

	if err != nil {
		exit(fmt.Sprintf("Something went wrong: %s", err.Error()))
	}

	correctAns := 0

	timerObj := time.NewTimer(time.Duration(timer) * time.Second)

	ans := make(chan string)
problemLoop:
	for i, p := range problems {
		var answer string
		fmt.Printf("Problem %d: %s=", i+1, p.question)
		go func() {
			fmt.Scanf("%s", &answer)
			ans <- answer
		}()

		select {
		case <-timerObj.C:
			fmt.Println()
			break problemLoop
		case iAns := <-ans:
			if iAns == p.answer {
				correctAns++
			}
			if i == len(problems)-1 {
				close(ans)
			}
		}
	}
	fmt.Printf("Your Result is %d out of %d\n", correctAns, len(problems))
	fmt.Printf("Press enter to exit")
	<-ans
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
