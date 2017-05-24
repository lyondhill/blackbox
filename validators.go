package blackbox

import (
	"strings"
	"fmt"
)

type validator func(stdout string, stderr string, success bool) error

func ValidStdout(str string) validator {
	return func(out, err string, success bool) error {
		if strings.Contains(out, str) {
			return nil
		}
		return fmt.Errorf("unable to location '%s' in stdout", str)
	}	
}

func ValidStderr(str string) validator {
	return func(out, err string, success bool) error {
		if strings.Contains(err, str) {
			return nil
		}
		return fmt.Errorf("unable to location '%s' in stderr", str)
	}	
}

func ValidOutput(str string) validator {
	return func(out, err string, success bool) error {
		if strings.Contains(out, str) || strings.Contains(err, str) {
			return nil
		}
		return fmt.Errorf("unable to location '%s' in output", str)
	}
}

func ValidExit() validator {
	return func(out, err string, success bool) error {
		if success {
			return nil		
		}
		return fmt.Errorf("unsuccessful exit")
	}	
}

func validNotOutput(str string) validator {
	return func(out, err string, success bool) error {
		if strings.Contains(out, str) || strings.Contains(err, str) {
			return fmt.Errorf("found '%s' in output", str)
		}
		return nil
	}	
}