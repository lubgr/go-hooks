package main

import (
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
)

type countingFwdWriter struct {
	io.Writer
	count int
}

// Write accumulates the received byte count and forwards to the instance's Writer.
func (w *countingFwdWriter) Write(p []byte) (n int, err error) {
	w.count += len(p)
	return w.Writer.Write(p)
}

// This program executes a given other program and forwards its stdout/stderr to its own
// stdout/stderr. When the subprocess terminates with a failure, that exit code is mirrored.
// Otherwise, when the subprocess printed anything to stdout or stderr, a failure exit code is used.
// When the subprocess terminates successfully and does not print anything to stdout/stderr, this
// program terminates successfully.
func main() {
	log.SetFlags(0)

	// No need to use actual argument passing with flags. The only thing we do is forwarding
	// everything to the executable given as the first argument.
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %v LINTER [flags...]", os.Args[0])
	}

	linter := exec.Command(os.Args[1], os.Args[2:]...)
	stdout := &countingFwdWriter{Writer: os.Stdout}
	stderr := &countingFwdWriter{Writer: os.Stderr}
	linter.Stdout = stdout
	linter.Stderr = stderr

	errRun := linter.Run()

	var errExit *exec.ExitError
	if errors.As(errRun, &errExit) {
		os.Exit(errExit.ExitCode())
	}

	if stdout.count > 0 || stderr.count > 0 {
		os.Exit(1)
	}
}
