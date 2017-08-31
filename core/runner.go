package core

import (
	"fmt"
	"github.com/gosuri/uilive"
	"time"
)

type TaskResult struct {
	task   Task
	result CmdResult
}

type Runner struct {
	config        Config
	fix           bool
	resultChannel chan TaskResult
	doneChannel   chan map[Task]TaskResult
}

func NewRunner(config Config, fix bool) *Runner {
	return &Runner{
		config:        config,
		fix:           fix,
		resultChannel: make(chan TaskResult),
		doneChannel:   make(chan map[Task]TaskResult),
	}
}

func (this Runner) Run() {
	for i := 0; i < len(this.config.Tasks); i += 1 {
		go this.processTask(this.config.Tasks[i])
	}
	go this.reportProgress()

	finalResults := <-this.doneChannel
	// Report
	for _, taskResult := range finalResults {
		fmt.Println("\nResults for", taskResult.task.Name)
		fmt.Println("success? %v", taskResult.result.success)
		fmt.Println(taskResult.result.output)
	}
}

func (this Runner) processTask(task Task) {
	// Use the FixCommand or regular Command depending on the flag passed to CLI
	cmdStr := task.Command

	if this.fix && task.FixCommand != "" {
		cmdStr = task.FixCommand
	}

	// Execute command
	result := NewCmd(cmdStr).Execute()

	this.resultChannel <- TaskResult{
		task:   task,
		result: result,
	}
}

func (this Runner) reportProgress() {
	writer := uilive.New()
	writer.Start()

	// Define a map of task_name => TaskResult
	results := make(map[Task]TaskResult)

	// Use a ticker here
	ticker := time.NewTicker(time.Millisecond * 50)

	i := 0
	for range ticker.C {
		// Check if there is a message on the channel
		select {
		case result := <-this.resultChannel:
			results[result.task] = result
		default:
		}

		// if so, update the status, if not still pending
		fmt.Fprintf(writer, this.generateProgressString(i, results))
		i += 1
		if len(results) == len(this.config.Tasks) {
			writer.Stop()

			this.doneChannel <- results
			return
		}
	}
}

func (this Runner) generateProgressString(tick int, results map[Task]TaskResult) string {
	var str = ""
	for i := 0; i < len(this.config.Tasks); i += 1 {
		task := this.config.Tasks[i]
		charSet := []string{"🕐", "🕑", "🕒", "🕓", "🕔", "🕕", "🕖", "🕗", "🕘", "🕙", "🕚", "🕛"}
		status := charSet[tick%len(charSet)]
		if result, ok := results[task]; ok {
			if result.result.success {
				status = "✅"
			} else {
				status = "❌"
			}
		}

		str += status + "  " + task.Name + "\n"
	}

	return str
}
