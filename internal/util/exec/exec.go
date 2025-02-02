package exec

import (
	"bytes"
	"encoding/json"
	"log/slog"

	execPackage "k8s.io/utils/exec"
)

type ExecInterface execPackage.Interface
type ExecCommand execPackage.Cmd

var exec execPackage.Interface

func init() {
	exec = execPackage.New()
}

func getExec() execPackage.Interface {
	return exec
}

func SetExec(e execPackage.Interface) {
	exec = e
}

type Command struct {
	ExecCommand
}

func NewCommand(bin string, args ...string) Command {
	cmd := getExec().Command(bin, args...)
	return Command{ExecCommand: cmd}
}

func (c Command) Run(stdin *bytes.Buffer, stdout *bytes.Buffer) error {
	if stdin != nil {
		c.SetStdin(stdin)
	}
	if stdout != nil {
		c.SetStdout(stdout)
	}
	stderr := bytes.NewBuffer(nil)
	c.SetStderr(stderr)
	err := c.ExecCommand.Run()
	stddErrStr := stderr.String()
	if stddErrStr != "" {
		slog.Info("command exec has stderr output", "error", stddErrStr)
	}

	return err
}

func (c Command) RunWithJSON(stdin interface{}, stdout interface{}) error {
	var stdinBuf *bytes.Buffer
	if stdin != nil {
		stdinStr, err := json.Marshal(stdin)
		if err != nil {
			return err
		}
		stdinBuf = bytes.NewBuffer(stdinStr)
	}
	stdoutBuf := bytes.NewBuffer(nil)
	err := c.Run(stdinBuf, stdoutBuf)
	if err != nil {
		return err
	}
	err = json.Unmarshal(stdoutBuf.Bytes(), stdout)
	if err != nil {
		return err
	}
	return nil
}
