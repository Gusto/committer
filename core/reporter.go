package core

import (
	"fmt"
	"github.com/gosuri/uilive"
	"github.com/kyokomi/emoji"
	"strings"
	"time"
)

type Reporter struct {
	tasks         []Task
	results       map[Task]TaskResult
	resultChannel chan TaskResult
	doneChannel   chan []TaskResult
}

func NewReporter(tasks []Task, resultChannel chan TaskResult) *Reporter {
	return &Reporter{
		tasks:         tasks,
		resultChannel: resultChannel,
		doneChannel:   make(chan []TaskResult),
	}
}

func (this Reporter) Report() {
	go this.reportProgress()
	this.reportFinalResults()
}

func (this Reporter) reportProgress() {
	writer := uilive.New()
	writer.Start()

	results := make(map[Task]TaskResult)

	// Use a ticker here
	ticker := time.NewTicker(time.Millisecond * 50)

	i := 0
	for range ticker.C {
		select {
		case result := <-this.resultChannel:
			results[result.task] = result
		default:
		}

		// If there is a message on the channel, pass along the result
		//    otherwise, continue to show the pending indicator
		emoji.Fprintf(writer, this.generateProgressString(i, results))

		if len(results) == len(this.tasks) {
			resultsArr := []TaskResult{}
			for _, taskResult := range results {
				resultsArr = append(resultsArr, taskResult)
			}
			writer.Stop()
			this.doneChannel <- resultsArr
			return
		}

		i += 1
	}
}

func (this Reporter) reportFinalResults() {
	finalResults := <-this.doneChannel

	var failures bool
	for _, taskResult := range finalResults {
		if !taskResult.success {
			if !failures {
				failures = true
			}
			fmt.Println("\nResults for", taskResult.task.Name)
			fmt.Println(taskResult.output + "\n")
		}

		if len(taskResult.fixedOutput) > 0 {
			fmt.Println("Autocorrected: ")
			fmt.Println(strings.TrimSpace(taskResult.fixedOutput) + "\n")
		}
	}

	if failures {
		emoji.Println(":x: Finished pre-commit hook")
	} else {
		emoji.Println(":white_check_mark: Finished pre-commit hook")
	}
}

func (this Reporter) generateProgressString(tick int, results map[Task]TaskResult) string {
	var str = ""
	for i := 0; i < len(this.tasks); i += 1 {
		task := this.tasks[i]
		clockSet := []string{
			":clock1:",
			":clock2:",
			":clock3:",
			":clock4:",
			":clock5:",
			":clock6:",
			":clock7:",
			":clock8:",
			":clock9:",
			":clock10:",
			":clock11:",
			":clock12:",
		}
		status := clockSet[tick%len(clockSet)]
		if result, ok := results[task]; ok {
			if result.success {
				status = ":white_check_mark:"
			} else {
				status = ":x:"
			}
		}

		str += "  " + status + " - " + task.Name + "\n"
	}

	return str
}
