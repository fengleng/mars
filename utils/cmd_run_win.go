//go:build windows
// +build windows

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
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP}

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	if err = cmd.Start(); err != nil {
		log.Error(err)
		return
	}

	done := make(chan error)
	go func() error {
		log.Info("waiting sub-process complete")
		done <- cmd.Wait()
		log.Info("sub-process exited")
		return nil
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
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP}

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
