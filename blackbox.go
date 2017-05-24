package blackbox

import (
	"os/exec"
	"bytes"
	"io"
	"os"
	"fmt"
)

type validator func(stdout string, stderr string, success bool) error

type errorList struct {
	errors []error
}

type Command struct {
	execCmd *exec.Cmd
	validators []validator
}

func Cmd(name string, args ...string) *Command {
	return &Command{
		execCmd: exec.Command(name, args...),
		validators: []validator{},
	}
}

func (cmd *Command) Run() error {

	// setup the buffers so we can collect the output
	outBuf :=  &bytes.Buffer{}
	errBuf :=  &bytes.Buffer{}
	outMulti := io.MultiWriter(outBuf, os.Stdout)
	errMulti := io.MultiWriter(errBuf, os.Stderr)

	// setup the reader and write
	cmd.execCmd.Stdout = outMulti
	cmd.execCmd.Stderr = errMulti
	cmd.execCmd.Stdin  = os.Stdin

	// run the command and catch any execution errors
	err := cmd.execCmd.Run()	
	if err != nil {
		return fmt.Errorf("failed to exec: %s", err)
	}

	// create a list of errors that will display more clearly
	errors := errorList{
		errors: []error{},
	}

	outString := outBuf.String()
	errString := errBuf.String()
	success := cmd.execCmd.ProcessState.Success()

	for _, validator := range cmd.validators {
		err := validator(outString, errString, success)	
		if err != nil {
			errors.errors = append(errors.errors, err)
		}
	}

	if len(errors.errors) != 0 {
		return errors
	}

	return nil
}

func (cmd *Command) AddValidator(val validator) {
	cmd.validators = append(cmd.validators, val)
}

func (err errorList) Error() string {
	if len(err.errors) == 0 {
		return ""
	}

	str := "Validation Failed: \n"
	for _, er := range err.errors {
		str = fmt.Sprintf("%s  %s\n", str, er)
	}

	return str
}