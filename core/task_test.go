package core

import (
	"github.com/stretchr/testify/assert"
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
	task.Fix.Files = ".txt"

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
