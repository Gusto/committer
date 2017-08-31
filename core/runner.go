package core

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
	var tasksToRun []Task

	for i := 0; i < len(this.config.Tasks); i += 1 {
		task := this.config.Tasks[i]
		if task.shouldRun(this.changed) {
			tasksToRun = append(tasksToRun, task)
			go this.processTask(task)
		}
	}

	NewReporter(
		tasksToRun,
		this.resultChannel,
	).Report()
}

func (this Runner) processTask(task Task) {
	this.resultChannel <- task.Execute(this.changed, this.fix)
}
