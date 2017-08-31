package core

import "fmt"

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
	this.resultChannel <- task.Execute(this.fix, this.changed)
}
