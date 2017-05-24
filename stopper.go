package blackbox

import (
	"strings"
)

type stopper func(data string) bool

type stopWriter struct {
	stop stopper
	cmd *Command
}


func (sw stopWriter) Write(p []byte) (n int, err error) {
	if sw.stop(string(p)) {
		close(sw.cmd.done)
	}

	// we dont ever want the stop writer to fail
	return len(p), nil
}

func StopOnOutput(str string) stopper {
	return func(data string) bool {
		return strings.Contains(data, str)
	}
}