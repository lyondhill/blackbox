package blackbox

import (
	"strings"
	"os"
	"runtime"
)

type stopper func(data string) bool

type stopWriter struct {
	stop stopper
	cmd *Command
}


func (sw stopWriter) Write(p []byte) (n int, err error) {
	if sw.stop(string(p)) {
		if sw.cmd != nil && sw.cmd.execCmd != nil && sw.cmd.execCmd.Process != nil {
			// add some checks so we dont panic
			if runtime.GOOS == "windows" {
				pipe, err := sw.cmd.execCmd.StdinPipe()
				if err != nil {
					pipe.Write([]byte(`^C`))
				}
			} else {
				sw.cmd.execCmd.Process.Signal(os.Interrupt)
			}
		}
	}

	// we dont ever want the stop writer to fail
	return len(p), nil
}

func StopOnOutput(str string) stopper {
	return func(data string) bool {
		return strings.Contains(data, str)
	}
}