package core

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
)

func TestShouldRunWithRelevantFile(t *testing.T) {
	origChangedFiles := changedFilesList
	changedFilesList = []string{"one.rb"}
	defer func() { changedFilesList = origChangedFiles }()

	task := Task{
		Name:    "task",
		Command: "run-task",
		Files:   ".rb",
	}
	task.Fix.Command = "run-fix"

	assert.True(t, task.shouldRun(), "It should run with relevant changed files")
}

func TestPrepareCommandFix(t *testing.T) {
	var task Task
	task.Fix.Command = "run-fix"

	assert.Equal(
		t,
		[]string{"run-fix"},
		task.prepareCommand(true),
		"It runs the fix command when fix is true",
	)
}

func TestPrepareCommand(t *testing.T) {
	origChangedFiles := changedFilesList
	changedFilesList = []string{"one.rb", "two.js", "three.txt"}
	defer func() { changedFilesList = origChangedFiles }()

	var task Task
	task.Command = "run-task"
	task.Files = ".txt"

	assert.Equal(
		t,
		[]string{"run-task", "three.txt"},
		task.prepareCommand(false),
		"It correctly passes only the relevant files",
	)
}

func TestPrepareFixedOutput(t *testing.T) {
	var task Task
	task.Fix.Output = "Fixed:"

	output := `
Linted: 1
Fixed: 2
Linted: 3
Fixed: 4
`
	assert.Equal(
		t,
		`Fixed: 2
Fixed: 4`,
		task.prepareFixedOutput(output),
		"It returns the relevant output lines",
	)
}

func TestRelevantChangedFiles(t *testing.T) {
	var task Task
	task.Files = ".txt"

	files := []string{
		"one.rb",
		"two.js",
		"three.go",
		"four.txt",
	}

	assert.Equal(
		t,
		[]string{"four.txt"},
		task.relevantChangedFiles(files),
		"It returns the relevant output lines",
	)
}

/*
	task.Execute tests
*/

func stubExecCommand(output []byte, success bool) {
	execCommand = func(command string, args ...string) ([]byte, error) {
		var error error
		if !success {
			error = errors.New("Boom")
		}
		return output, error
	}
}

func restoreExecCommand() {
	execCommand = func(command string, args ...string) ([]byte, error) {
		return exec.Command(command, args...).CombinedOutput()
	}
}

func TestExecuteSuccess(t *testing.T) {
	stubExecCommand([]byte("Output!"), true)
	defer restoreExecCommand()

	task := Task{
		Name:    "task",
		Command: "run-task",
	}

	result := task.Execute(false)
	assert.True(t, result.success, "The result is successful")
	assert.Equal(t, result.task, task, "It attaches the task")
	assert.Equal(t, result.output, "Output!", "It attaches the task")
	assert.Equal(t, result.fixedOutput, "", "There is no fixed output")
}

func TestExecuteFailure(t *testing.T) {
	stubExecCommand([]byte("Output!"), false)
	defer restoreExecCommand()

	task := Task{
		Name:    "task",
		Command: "run-task",
	}

	result := task.Execute(false)
	assert.False(t, result.success, "The result is failed")
	assert.Equal(t, result.task, task, "It attaches the task")
	assert.Equal(t, result.output, "Output!", "It attaches the task")
	assert.Equal(t, result.fixedOutput, "", "There is no fixed output")
}

func TestExecuteFixSuccessNoFixCommand(t *testing.T) {
	stubExecCommand([]byte("Output!"), true)
	defer restoreExecCommand()

	task := Task{
		Name:    "task",
		Command: "run-task",
	}

	result := task.Execute(true)
	assert.True(t, result.success, "The result is successful")
	assert.Equal(t, result.task, task, "It attaches the task")
	assert.Equal(t, result.output, "Output!", "It does not grep through the output")
	assert.Equal(t, result.fixedOutput, "", "There is no fixed output")
}

func TestExecuteFixSuccessWithFixCommand(t *testing.T) {
	stubExecCommand(
		[]byte(`Linted: app/one.rb
Fixed: app/two.rb
Linted: app/three.rb
`),
		true,
	)
	defer restoreExecCommand()

	task := Task{
		Name:    "task",
		Command: "run-task",
	}
	task.Fix.Command = "run-fix"
	task.Fix.Output = "Fixed:"

	result := task.Execute(true)
	assert.True(t, result.success, "The result is successful")
	assert.Equal(t, result.task, task, "It attaches the task")
	assert.Equal(t, result.output, "Linted: app/one.rb\nFixed: app/two.rb\nLinted: app/three.rb\n", "It attaches the entire output")
	assert.Equal(t, result.fixedOutput, "Fixed: app/two.rb", "There is a subset of the output")
}

func TestExecuteFixSuccessWithFixCommandWithNoOuput(t *testing.T) {
	stubExecCommand([]byte(""), true)
	defer restoreExecCommand()

	task := Task{
		Name:    "task",
		Command: "run-task",
	}
	task.Fix.Command = "run-fix"
	task.Fix.Output = "Fixed:"

	result := task.Execute(true)
	assert.True(t, result.success, "The result is successful")
	assert.Equal(t, result.task, task, "It attaches the task")
	assert.Equal(t, result.output, "", "There is no output")
	assert.Equal(t, result.fixedOutput, "", "There is no output")
}

func TestExecuteFixSuccessWithFixCommandWithNoOuputWithAutostage(t *testing.T) {
	stubExecCommand([]byte(""), true)
	defer restoreExecCommand()

	task := Task{
		Name:    "task",
		Command: "run-task",
	}
	task.Fix.Command = "run-fix"
	task.Fix.Output = "Fixed:"
	task.Fix.Autostage = true

	result := task.Execute(true)
	assert.True(t, result.success, "The result is successful")
	assert.Equal(t, result.task, task, "It attaches the task")
	assert.Equal(t, result.output, "", "There is no output")
	assert.Equal(t, result.fixedOutput, "No output but staging since autostage is true", "There is a subset of the output")
}

func TestExecuteFixFailureWithFixCommand(t *testing.T) {
	stubExecCommand(
		[]byte("Failed!"),
		false,
	)
	defer restoreExecCommand()

	task := Task{
		Name:    "task",
		Command: "run-task",
	}
	task.Fix.Command = "run-fix"
	task.Fix.Output = "Fixed:"

	result := task.Execute(true)
	assert.False(t, result.success, "The result is not successful")
	assert.Equal(t, result.task, task, "It attaches the task")
	assert.Equal(t, result.output, "Failed!", "It attaches the entire output")
	assert.Equal(t, result.fixedOutput, "", "There is no fixed output")
}

func TestExecuteFixSuccessWithFixCommandShouldNotStage(t *testing.T) {
	origShouldStage := shouldStage
	shouldStage = false
	defer func() { shouldStage = origShouldStage }()

	stubExecCommand(
		[]byte(`Linted: app/one.rb
Fixed: app/two.rb
Linted: app/three.rb
`),
		true,
	)
	defer restoreExecCommand()

	task := Task{
		Name:    "task",
		Command: "run-task",
	}
	task.Fix.Command = "run-fix"
	task.Fix.Output = "Fixed:"

	result := task.Execute(true)
	assert.False(t, result.success, "The result is marked unsuccessful so changes can be staged")
	assert.Equal(t, result.task, task, "It attaches the task")
	assert.Equal(t, result.output, "Linted: app/one.rb\nFixed: app/two.rb\nLinted: app/three.rb\n", "It attaches the entire output")
	assert.Equal(t, result.fixedOutput, "Fixed: app/two.rb", "There is a subset of the output")
}
