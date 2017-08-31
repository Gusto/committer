package core

import "fmt"

type TaskResult struct {
	task   Task
	result CmdResult
}

type Runner struct {
	config        Config
	fix           bool
	changed       bool
	resultChannel chan TaskResult
}

func NewRunner(config Config, fix bool, changed bool) *Runner {
	return &Runner{
		config:        config,
		fix:           fix,
		changed:       changed,
		resultChannel: make(chan TaskResult),
	}
}

func (this Runner) Run() {
	fmt.Println("Running commit hook for:")

	for i := 0; i < len(this.config.Tasks); i += 1 {
		go this.processTask(this.config.Tasks[i])
	}

	NewReporter(
		this.config.Tasks,
		this.resultChannel,
	).Report()
}

func (this Runner) processTask(task Task) {
	// Use the FixCommand or regular Command depending on the flag passed to CLI
	var cmdStr string
	if this.fix && task.FixCommand != "" {
		cmdStr = task.FixCommand
	} else {
		cmdStr = task.Command
	}

	// Feed in changed files if we are running with --changed
	if this.changed {
		cmdStr += " " + NewCmd("git diff --cached --name-only").Execute().output
	}

	// Execute command
	result := NewCmd(cmdStr).Execute()

	// Handle autocorrect parsing here
	if this.fix && result.success {
		// if len(result.fixedFiles) > 0 {
		//
		// }
	}

	this.resultChannel <- TaskResult{
		task:   task,
		result: result,
	}
}
