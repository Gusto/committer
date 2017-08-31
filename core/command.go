package core

import (
	"os/exec"
	"strings"
)

type Cmd struct {
	Command string
}

type CmdResult struct {
	cmd     Cmd
	output  string
	success bool
}

func NewCmd(command string) *Cmd {
	return &Cmd{
		Command: command,
	}
}

func (cmd Cmd) Execute() *CmdResult {
	command := strings.Split(cmd.Command, " ")
	exeCommand := exec.Command(command[0], command[1:]...)

	output, err := exeCommand.CombinedOutput()

	// Return a CmdResult object
	return &CmdResult{
		cmd:     cmd,
		output:  string(output),
		success: err == nil,
	}
}
