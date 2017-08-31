package core

import (
	"os/exec"
	"strings"
)

type Task struct {
	Name       string `yaml:"name"`
	Command    string `yaml:"command"`
	FixCommand string `yaml:"fix_command"`
	FixGrep    string `yaml:"fix_grep"`
}

type TaskResult struct {
	task        Task
	success     bool
	output      string
	fixedOutput string
}

const shouldStage = true

func (task Task) Execute(fix bool, changed bool) TaskResult {
	// Use the FixCommand or regular Command depending on the flag passed to CLI
	var cmdStr string
	if fix && task.FixCommand != "" {
		cmdStr = task.FixCommand
	} else {
		cmdStr = task.Command
	}

	// Feed in changed files if we are running with --changed
	changedFiles, err := exec.Command("git", "diff", "--cached", "--name-only").Output()
	changedStr := strings.TrimSpace(string(changedFiles))
	if changed {
		if err != nil {
			panic(err)
		}
		cmdStr += " -- " + changedStr
	}

	// Execute command
	command := strings.Split(cmdStr, " ")
	exeCommand := exec.Command(command[0], command[1:]...)
	output, err := exeCommand.CombinedOutput()

	outputStr := string(output)
	success := err == nil

	// Handle autocorrect parsing here
	var fixedOutputList []string
	if fix && success {
		for _, item := range strings.Split(outputStr, "\n") {
			if strings.Contains(item, task.FixGrep) {
				fixedOutputList = append(fixedOutputList, item)
			}
		}

		if len(fixedOutputList) > 0 {
			if shouldStage {
				changedFilesList := strings.Split(changedStr, " ")
				subCmd := append([]string{"add"}, changedFilesList...)

				if _, err := exec.Command("git", subCmd...).Output(); err != nil {
					panic(err)
				}
			} else {

			}
		}
	}

	// Return a CmdResult object
	return TaskResult{
		task:        task,
		output:      outputStr,
		success:     success,
		fixedOutput: strings.Join(fixedOutputList, "\n"),
	}
}
