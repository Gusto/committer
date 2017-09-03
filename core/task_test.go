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

	assert.True(t, task.shouldRun(false), "It should run when changed")
}

func TestShouldRunNotChangedNoFix(t *testing.T) {
	task := Task{
		Name:    "task",
		Command: "run-task",
	}

	assert.False(t, task.shouldRun(true), "Its true")
}
