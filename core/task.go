package core

import (
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type TaskFix struct {
}
type Task struct {
	Name    string
	Command string
	Fix     struct {
		Command string
		Output  string
		Files   string
	}
}

type TaskResult struct {
	task        Task
	success     bool
	output      string
	fixedOutput string
}

var shouldStage = (os.Getenv("COMMITTER_SKIP_STAGE_FIX") == "")
var changedFiles, _ = exec.Command("git", "diff", "--cached", "--name-only").Output()
var changedFilesList = strings.Split(string(changedFiles), "\n")

func (task Task) Execute(changed bool, fix bool) TaskResult {
	// Use the FixCommand or regular Command depending on the flag passed to CLI
	var cmdStr string
	if fix && task.Fix.Command != "" {
		cmdStr = task.Fix.Command
	} else {
		cmdStr = task.Command
	}

	// Feed in changed files if we are running with --changed
	var changedStr string
	var relevantChangedFilesList []string

	if changed {
		for _, file := range changedFilesList {
			match, _ := regexp.MatchString(task.Fix.Files, file)
			if match {
				relevantChangedFilesList = append(relevantChangedFilesList, file)
			}
		}

		changedStr = strings.Join(relevantChangedFilesList, " ")
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

	Strings:
		for _, item := range strings.Split(outputStr, "\n") {
			match, _ := regexp.MatchString(task.Fix.Output, item)

			if match {
				fixedOutputList = append(fixedOutputList, item)
				continue Strings
			}
		}

		if len(fixedOutputList) > 0 {
			if shouldStage {
				subCmd := append([]string{"add"}, relevantChangedFilesList...)

				if _, err := exec.Command("git", subCmd...).Output(); err != nil {
					panic(err)
				}
			} else {
				// Explicitly mark autocorrects that were not stage as unsuccessful
				//   so they can be staged manually
				success = false
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

func (task Task) shouldRun(changed bool) bool {
	// Always run all tasks if we aren't just looking at changed files
	if !changed {
		return true
	}

	if task.Fix.Command != "" {
		for _, file := range changedFilesList {
			match, err := regexp.MatchString(task.Fix.Files, file)

			if err != nil {
				panic(err)
			}
			if match {
				return true
			}
		}
	}

	return false
}
