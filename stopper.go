package blackbox

import (
	"strings"
	"os"
)

type stopper func(data string) bool

type stopWriter struct {
	stop stopper
	cmd *Command
}


// func terminateProcess(pid, exitcode int) error {
// 	h, e := syscall.OpenProcess(syscall.PROCESS_TERMINATE, false, uint32(pid))
// 	if e != nil {
// 		return NewSyscallError("OpenProcess", e)
// 	}
// 	defer syscall.CloseHandle(h)
// 	e = syscall.TerminateProcess(h, uint32(exitcode))
// 	return NewSyscallError("TerminateProcess", e)
// }

func (sw stopWriter) Write(p []byte) (n int, err error) {
	if sw.stop(string(p)) {
		if sw.cmd != nil && sw.cmd.execCmd != nil && sw.cmd.execCmd.Process != nil {
			// add some checks so we dont panic

			killChild(Process)
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