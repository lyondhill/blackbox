package blackbox

import (
	"os/exec"
	"bytes"
	"io"
	"os"
	"fmt"
)

var Quiet bool


type errorList struct {
	errors []error
}

type Command struct {
	execCmd *exec.Cmd
	validators []validator
	stoppers []stopper
	done chan struct{}
}

func Cmd(name string, args ...string) *Command {
	return &Command{
		execCmd: exec.Command(name, args...),
		validators: []validator{},
		stoppers: []stopper{},
		done: make(chan struct{}, 1),
	}
}

func (cmd *Command) Run() error {

	outBuf :=  &bytes.Buffer{}
	errBuf :=  &bytes.Buffer{}

	outWriter := []io.Writer{outBuf}
	errwriter := []io.Writer{errBuf}

	for _, stopper := range cmd.stoppers {
		sw := stopWriter{
			cmd: cmd,
			stop: stopper,
		}
		outWriter = append(outWriter, sw)
		errwriter = append(errwriter, sw)
	}

	if !Quiet {
		outWriter = append(outWriter, os.Stdout)
		errwriter = append(errwriter, os.Stderr)
	}

	// setup the buffers so we can collect the output
	outMulti := io.MultiWriter(outWriter...)
	errMulti := io.MultiWriter(errwriter...)

	// setup the reader and write
	cmd.execCmd.Stdout = outMulti
	cmd.execCmd.Stderr = errMulti
	cmd.execCmd.Stdin  = os.Stdin

	// run the command and catch any execution errors
	cmd.execCmd.Run()	
	// if err != nil {
	// 	return fmt.Errorf("failed to exec: %s", err)
	// }

	// execDone := make(chan error, 1)
	// go func() {
 //    	execDone <- cmd.execCmd.Wait()
	// }()

	// for {
	// 	select {
	// 	case <-cmd.done:
	// 		// a stopper triggered a stop here
	// 		// kill the command
	// 		cmd.execCmd.Process.Kill()
	// 	case err = <-execDone:
	// 		break
	// 	}		
	// }

	// if err != nil {
	// 	return fmt.Errorf("failed on wait: %s", err)
	// }

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

func (cmd *Command) AddStopper(stop stopper) {
	cmd.stoppers = append(cmd.stoppers, stop)
	
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