// +build !windows

package blackbox

import "os"

func killChild(p os.Process) error {
	return p.Signal(os.Interrupt)
}
