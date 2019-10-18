package core

type Runner struct {
	config        Config
	fix           bool
	resultChannel chan TaskResult
}

func NewRunner(config Config, fix bool) *Runner {
	return &Runner{
		config:        config,
		fix:           fix,
		resultChannel: make(chan TaskResult),
	}
}

func (this Runner) Run() bool {
	var tasksToRun []Task

	for i := 0; i < len(this.config.Tasks); i += 1 {
		task := this.config.Tasks[i]
		if task.shouldRun() {
			tasksToRun = append(tasksToRun, task)
			go this.processTask(task)
		}
	}

	success := NewReporter(
		tasksToRun,
		this.resultChannel,
	).Report()

	return success
}

func (this Runner) processTask(task Task) {
	this.resultChannel <- task.Execute(this.fix)
}
