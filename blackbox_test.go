package blackbox

import (
	"testing"
)

func TestRun(t *testing.T) {
	cmd := Cmd("echo","hello")
	cmd.AddValidator(ValidExit())
	cmd.AddValidator(ValidStdout("hello"))
	err := cmd.Run()
	if err != nil {
		t.Error(err)
	}
}

func TestBadRun(t *testing.T) {
	cmd := Cmd("echo", "howdy")
	cmd.AddValidator(ValidOutput("hello"))
	err := cmd.Run()
	if err == nil {
		t.Errorf("should have failed because of stdout mismatch")
	}
	
}

func TestStopper(t *testing.T) {
	cmd := Cmd("top", "-b")
	cmd.AddStopper(StopOnOutput("Tasks:"))
	cmd.AddValidator(ValidExit())
	err := cmd.Run()
	if err != nil {
		t.Error(err)
	}
}