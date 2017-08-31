package core

import (
	"fmt"
	"github.com/gosuri/uilive"
	"sync"
	"time"
)

type TaskResult struct {
	task   *Task
	result *CmdResult
}

func Run(config *Config, fix bool) {
	var wg sync.WaitGroup

	var resultsMutex = &sync.Mutex{}
	var results []TaskResult

	var doneChannel = make(chan TaskResult)

	wg.Add(1)
	go generateProgress(config.Tasks, doneChannel, &wg)

	for i := 0; i < len(config.Tasks); i += 1 {
		// copyi := i
		task := config.Tasks[i]
		wg.Add(1)
		go func() {
			// Use the FixCommand or regular Command depending on the flag passed to CLI
			cmdStr := task.Command
			if fix && task.FixCommand != "" {
				cmdStr = task.FixCommand
			}

			// Execute command
			result := NewCmd(cmdStr).Execute()

			// Update the results array
			resultsMutex.Lock()
			defer resultsMutex.Unlock()
			// fmt.Println("Appending: %v", task.Name)
			taskResult := &TaskResult{
				task:   &task,
				result: result,
			}
			results = append(results, *taskResult)
			doneChannel <- *taskResult
			// Mark the action complete
			wg.Done()
		}()
	}

	// Wait for all commands to complete
	wg.Wait()

	// Report
	for i := 0; i < len(results); i += 1 {
		taskResult := results[i]

		fmt.Println("\nResults for", taskResult.task.Name)
		fmt.Println("success? %v", taskResult.result.success)
		fmt.Println(taskResult.result.output)
	}
}

func generateProgress(tasks []Task, doneChannel chan TaskResult, wg *sync.WaitGroup) {
	writer := uilive.New()
	writer.Start()

	// Define a map of task_name => TaskResult
	results := make(map[Task]TaskResult)

	// Use a ticker here
	ticker := time.NewTicker(time.Millisecond * 500)

	for range ticker.C {
		// Check if there is a message on the channel
		select {
		case result := <-doneChannel:
			// fmt.Println("received message", result)
			results[*result.task] = result
		default:
			// fmt.Println("no message received")
		}

		var str = ""
		for i := 0; i < len(tasks); i += 1 {
			task := tasks[i]
			status := "" // / - \ - /
			if result, ok := results[task]; ok {
				if result.result.success {
					status = "Completed"
				} else {
					status = "Failed"
				}
			}
			str += task.Name + "...... " + status + "\n"
		}
		// if so, update the status, if not still pending
		fmt.Fprintf(writer, str)

		if len(results) == len(tasks) {
			writer.Stop()
			wg.Done()
			// close(doneChannel)
			return
		}
	}
}
