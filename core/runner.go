package core

import (
	"fmt"
	"sync"
)

type TaskResult struct {
	task   Task
	result *CmdResult
}

func Run(config *Config, fix bool) {
	var wg sync.WaitGroup

	var resultsMutex = &sync.Mutex{}
	var results []*TaskResult

	for i := 0; i < len(config.Tasks); i += 1 {
		task := config.Tasks[i]
		fmt.Println(task.Name)

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
			fmt.Println("Appending: %v", task.Name)
			results = append(results, &TaskResult{
				task:   task,
				result: result,
			})

			// Mark the action complete
			wg.Done()
		}()
	}

	// Wait for all commands to complete
	wg.Wait()

	// Rreport
	for i := 0; i < len(results); i += 1 {
		taskResult := results[i]

		fmt.Println("\nResults for", taskResult.task.Name)
		fmt.Println("success? %v", taskResult.result.success)
		fmt.Println(taskResult.result.output)
	}
}

func reportOutput(results []*TaskResult) {

}
