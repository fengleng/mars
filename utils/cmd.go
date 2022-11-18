package utils

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gososy/sorpc/log"
)

var ErrRunCommandTimeout = errors.New("run command timeout")

func RunCommandWithTimeout(timeout time.Duration, name string, args ...string) (err error, stdout string, stderr string, exitStatus int) {
	cmd := exec.Command(name, args...)
	cmd.Env = os.Environ()

	return runCommandWithTimeOut(cmd, timeout)
}

func RunCommandWithTimeoutWorkDir(timeout time.Duration, workDir, name string, args ...string) (err error, stdout string, stderr string, exitStatus int) {
	cmd := exec.Command(name, args...)
	cmd.Dir = workDir
	cmd.Env = os.Environ()

	return runCommandWithTimeOut(cmd, timeout)
}

func RunCommandWithWorkDir(workDir, name string, args ...string) (err error, stdout string, stderr string, exitStatus int) {
	cmd := exec.Command(name, args...)
	cmd.Dir = workDir
	cmd.Env = os.Environ()
	log.Info(strings.Join(cmd.Args, " "))
	log.Infof("run at:%s", cmd.Dir)

	return runCommand(cmd)
}

func RunCommand(name string, args ...string) (err error, stdout string, stderr string, exitStatus int) {
	cmd := exec.Command(name, args...)
	cmd.Env = os.Environ()
	log.Info(strings.Join(cmd.Args, " "))

	return runCommand(cmd)
}
