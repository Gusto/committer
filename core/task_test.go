package core

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
)

func TestShouldRunNotChanged(t *testing.T) {
	task := Task{
		Name:    "task",
		Command: "run-task",
	}

	assert.True(t, task.shouldRun(false), "It runs when not in changed mode")
}

func TestShouldRunNotChangedNoFix(t *testing.T) {
	task := Task{
		Name:    "task",
		Command: "run-task",
	}

	assert.False(t, task.shouldRun(true), "It does not run when in changed mode with no fix command")
}

func TestShouldRunNotChangedWithFix(t *testing.T) {
	task := Task{
		Name:    "task",
		Command: "run-task",
	}
	task.Fix.Command = "run-fix"

	assert.True(t, task.shouldRun(true), "It does not run when in changed mode with no fix command")
}

func TestPrepareCommandNoChangeNoFix(t *testing.T) {
	task := Task{
		Command: "run-task",
	}

	assert.Equal(
		t,
		[]string{"run-task"},
		task.prepareCommand(false, true),
		"It runs the fix command when fix is true",
	)
}

func TestPrepareCommandNoChangeFix(t *testing.T) {
	var task Task
	task.Fix.Command = "run-fix"

	assert.Equal(
		t,
		[]string{"run-fix"},
		task.prepareCommand(false, true),
		"It runs the fix command when fix is true",
	)
}

func TestPrepareCommandWithChanged(t *testing.T) {
	origChangedFiles := changedFilesList
	changedFilesList = []string{"one.rb", "two.js", "three.txt"}
	defer func() { changedFilesList = origChangedFiles }()

	var task Task
	task.Command = "run-task"
	task.Files = ".txt"

	assert.Equal(
		t,
		[]string{"run-task", "--", "three.txt"},
		task.prepareCommand(true, false),
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
type Executor func(command string, args ...string) ([]byte, error)

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

	result := task.Execute(false, false)
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

	result := task.Execute(false, false)
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

	result := task.Execute(false, true)
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

	result := task.Execute(false, true)
	assert.True(t, result.success, "The result is successful")
	assert.Equal(t, result.task, task, "It attaches the task")
	assert.Equal(t, result.output, "Linted: app/one.rb\nFixed: app/two.rb\nLinted: app/three.rb\n", "It attaches the entire output")
	assert.Equal(t, result.fixedOutput, "Fixed: app/two.rb", "There is a subset of the output")
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

	result := task.Execute(false, true)
	assert.False(t, result.success, "The result is successful")
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

	result := task.Execute(false, true)
	assert.False(t, result.success, "The result is marked unsuccessful so changes can be staged")
	assert.Equal(t, result.task, task, "It attaches the task")
	assert.Equal(t, result.output, "Linted: app/one.rb\nFixed: app/two.rb\nLinted: app/three.rb\n", "It attaches the entire output")
	assert.Equal(t, result.fixedOutput, "Fixed: app/two.rb", "There is a subset of the output")
}
