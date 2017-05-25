package blackbox

import  (
	"syscall"
	"os"
	"fmt"
)

func killChild(p os.Process) error {
	var kernel32 = syscall.NewLazyDLL("Kernel32.dll")
	var freeConsole = kernel32.NewProc("FreeConsole")
	var attachConsole = kernel32.NewProc("AttachConsole")
	var setConsoleCtrlHandler = kernel32.NewProc("SetConsoleCtrlHandler")
	var generateConsoleCtrlEvent = kernel32.NewProc("GenerateConsoleCtrlEvent")
	// Close current console
	if r, _, err := freeConsole.Call(); r == 0 {
		return fmt.Errorf("Can't FreeConsole. Error code %v", err)
	}
	// Stach to job console
	if r, _, err := attachConsole.Call(uintptr(p.Pid)); r == 0 {
		return fmt.Errorf("Can't AttachConsole. Error code %v", err)
	}
	// Disable ctrl+C handling for our own program, so we don't "kill" ourselves
	if r, _, err := setConsoleCtrlHandler.Call(0, 1); r == 0 {
		return fmt.Errorf("Can't SetConsoleCtrlHandler. Error code %v", err)
	}

	if r, _, err := generateConsoleCtrlEvent.Call(uintptr(syscall.CTRL_C_EVENT), 0); r == 0 {
		return fmt.Errorf("Can't GenerateConsoleCtrlEvent. Error code %v", err)
	}
	return nil	
}
