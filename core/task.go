package core

import (
	"os"
	"os/exec"
	"strings"
)

type Task struct {
	Name           string   `yaml:"name"`
	Command        string   `yaml:"command"`
	FixCommand     string   `yaml:"fix_command"`
	FixGrep        []string `yaml:"fix_grep"`
	FileExtensions []string `yaml:"file_extensions"`
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
	if fix && task.FixCommand != "" {
		cmdStr = task.FixCommand
	} else {
		cmdStr = task.Command
	}

	// Feed in changed files if we are running with --changed
	changedStr := strings.Join(changedFilesList, " ")
	if changed {
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
			for _, keyword := range task.FixGrep {
				if strings.Contains(item, keyword) {
					fixedOutputList = append(fixedOutputList, item)
					continue Strings
				}
			}
		}

		if len(fixedOutputList) > 0 {
			if shouldStage {
				changedFilesList := strings.Split(changedStr, " ")
				changedFilesList = changedFilesList[:(len(changedFilesList) - 1)]
				subCmd := append([]string{"add"}, changedFilesList...)

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

	for _, file := range changedFilesList {
		for _, suffix := range task.FileExtensions {
			if strings.HasSuffix(file, suffix) {
				return true
			}
		}
	}

	return false
}
