package core

import (
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type Task struct {
	Name             string
	Command          string
	Files            string
	ExcludeFilenames bool
	Fix              struct {
		Command   string
		Output    string
		Autostage bool
	}
}

type TaskResult struct {
	task        Task
	success     bool
	output      string
	fixedOutput string
}

var shouldStage = (os.Getenv("COMMITTER_SKIP_STAGE") == "")
var autoFix = (os.Getenv("COMMITTER_AUTO_FIX") == "")

var changedFiles, _ = exec.Command("git", "diff", "--diff-filter=ACMRTUXB", "--cached", "--name-only").Output()
var changedFilesList = strings.Split(string(changedFiles), "\n")

func (task Task) relevantChangedFiles(changedFilesList []string) []string {
	var relevantChangedFilesList []string

	for _, file := range changedFilesList {
		match, _ := regexp.MatchString(task.Files, file)
		if match && file != "" {
			relevantChangedFilesList = append(relevantChangedFilesList, file)
		}
	}

	return relevantChangedFilesList
}

var execCommand = func(command string, args ...string) ([]byte, error) {
	return exec.Command(command, args...).CombinedOutput()
}

func (task Task) Execute(fix bool) TaskResult {
	// Generate command based on --fix / --changed
	autoFixORFix := autoFix || fix
	command := task.prepareCommand(autoFixORFix)

	// Run command
	output, err := execCommand(command[0], command[1:]...)

	outputStr := string(output)
	success := err == nil

	var fixedOutput string
	shouldFix := autoFixORFix && task.Fix.Command != ""
	if success && shouldFix {
		// If we are fixing and successfully updated files, capture the output
		fixedOutput = task.prepareFixedOutput(outputStr)
		if fixedOutput == "" && task.Fix.Autostage {
			fixedOutput = "No output but staging since autostage is true"
		}

		if fixedOutput != "" {
			// If we have output, then we've corrected files
			if shouldStage {
				// Stage files by default automatically
				task.stageRelevantFiles()
			} else {
				// Explicitly fail the pre-commit hook so the files can be staged manually
				success = false
			}
		}
	}

	return TaskResult{
		task:        task,
		success:     success,
		output:      outputStr,
		fixedOutput: fixedOutput,
	}
}

func (task Task) prepareFixedOutput(outputStr string) string {
	var fixedOutputList []string

	for _, item := range strings.Split(outputStr, "\n") {
		match, _ := regexp.MatchString(task.Fix.Output, item)

		if match {
			fixedOutputList = append(fixedOutputList, item)
			continue
		}
	}

	return strings.Join(fixedOutputList, "\n")
}

func (task Task) prepareCommand(fix bool) []string {
	// Use the FixCommand or regular Command depending on the flag passed to CLI
	var cmdStr string
	if fix && task.Fix.Command != "" {
		cmdStr = task.Fix.Command
	} else {
		cmdStr = task.Command
	}

	// Feed in changed files if we are running with --changed
	relevantChangedFilesList := task.relevantChangedFiles(changedFilesList)
	if !task.ExcludeFilenames && len(relevantChangedFilesList) > 0 {
		cmdStr += " " + strings.Join(relevantChangedFilesList, " ")
	}

	return strings.Split(cmdStr, " ")
}

func (task Task) stageRelevantFiles() {
	relevantChangedFiles := task.relevantChangedFiles(changedFilesList)
	subCmd := append([]string{"add"}, relevantChangedFiles...)

	if _, err := execCommand("git", subCmd...); err != nil {
		panic(err)
	}
}

func (task Task) shouldRun() bool {
	for _, file := range changedFilesList {
		match, err := regexp.MatchString(task.Files, file)

		if err != nil {
			panic(err)
		}
		if match {
			return true
		}
	}
	return false
}
