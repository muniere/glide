package shell

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

type CaptureResult struct {
	Success bool
	Code    int
	Stdout  string
	Stderr  string
}

type CallResult struct {
	Success bool
	Code    int
}

func Call(command string, args []string) CallResult {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err == nil {
		return CallResult{Success: true, Code: 0}
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return CallResult{Success: false, Code: exitErr.ExitCode()}
	}

	return CallResult{Success: false, Code: 1}
}

func Capture(command string, args ...string) CaptureResult {
	cmd := exec.Command(command, args...)
	out, err := cmd.CombinedOutput()
	text := strings.TrimRight(string(out), "\n")
	if err == nil {
		return CaptureResult{Success: true, Code: 0, Stdout: text, Stderr: ""}
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return CaptureResult{Success: false, Code: exitErr.ExitCode(), Stdout: "", Stderr: text}
	}

	return CaptureResult{Success: false, Code: 1, Stdout: "", Stderr: text}
}
