//go:build linux || darwin
// +build linux darwin

package utils

import (
	"bytes"
	"os/exec"
	"syscall"
	"time"

	"github.com/gososy/sorpc/log"
)

func runCommandWithTimeOut(cmd *exec.Cmd, timeout time.Duration) (err error, stdout string, stderr string, exitStatus int) {
	// https://stackoverflow.com/questions/392022/whats-the-best-way-to-send-a-signal-to-all-members-of-a-process-group
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	if err = cmd.Start(); err != nil {
		log.Error(err)
		return
	}

	done := make(chan error)
	go func() {
		log.Debug("waiting sub-process complete")
		done <- cmd.Wait()
		log.Debug("sub-process exited")
	}()

	isTimeout := false
	select {
	case err = <-done:
		// exited
		stdout = outBuf.String()
		stderr = errBuf.String()

		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
					exitStatus = status.ExitStatus()
				}
			}
		}

	case <-time.After(timeout):
		isTimeout = true
		if cmd.Process != nil {
			err = cmd.Process.Kill()
		} else {
			err = ErrRunCommandTimeout
		}
	}

	if isTimeout {
		<-done
	}

	return
}

func runCommand(cmd *exec.Cmd) (err error, stdout string, stderr string, exitStatus int) {
	// https://stackoverflow.com/questions/392022/whats-the-best-way-to-send-a-signal-to-all-members-of-a-process-group
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	if err = cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				exitStatus = status.ExitStatus()
			}
		}
	}
	stdout = outBuf.String()
	stderr = errBuf.String()
	return
}
