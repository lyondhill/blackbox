package blackbox

import (
	"os/exec"
)

type validator func(stdout string, stderr string, exit int) bool

type Command struct {
	execCmd *exec.Cmd
	validators []validator
}

func Cmd(name string, args ...string) *Command {
	return &Command{
		execCmd: exec.Cmd(name, args...),
		validators: []validator{},
	}
}


func

func (cmd *Command) AddValidator(val validator) {
	cmd.validators = append(validators, val)
}