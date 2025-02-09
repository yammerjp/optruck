package exec

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
	bin  string
	args []string
}

func NewCommand(bin string, args ...string) Command {
	cmd := getExec().Command(bin, args...)
	return Command{ExecCommand: cmd, bin: bin, args: args}
}

func sealJsonValues(stdin string) string {
	var data interface{}
	if err := json.Unmarshal([]byte(stdin), &data); err != nil {
		return stdin
	}

	sealedData := sealValues(data)
	sealed, err := json.Marshal(sealedData)
	if err != nil {
		return stdin
	}
	return string(sealed)
}

func sealValues(data interface{}) interface{} {
	switch v := data.(type) {
	case string:
		return "*****"
	case map[string]interface{}:
		sealed := make(map[string]interface{})
		for key, value := range v {
			sealed[key] = sealValues(value)
		}
		return sealed
	case []interface{}:
		sealed := make([]interface{}, len(v))
		for i, value := range v {
			sealed[i] = sealValues(value)
		}
		return sealed
	default:
		return v
	}
}

func (c Command) Run(stdin *bytes.Buffer, stdout *bytes.Buffer) error {
	if stdin != nil {
		// credentials are have to include json values for sealing
		slog.Debug("set stdin", "stdin", sealJsonValues(stdin.String()))
		c.SetStdin(stdin)
	}
	if stdout != nil {
		c.SetStdout(stdout)
	}
	stderr := bytes.NewBuffer(nil)
	c.SetStderr(stderr)
	slog.Info("run command", "bin", c.bin, "args", c.args)
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
		if errors.Is(err, execPackage.ErrExecutableNotFound) {
			return fmt.Errorf("command not found, please install the command `%s`: %w. Ensure the command is installed and try again.", c.bin, err)
		}
		return fmt.Errorf("failed to run command `%s`: %w. Please check the command and try again.", c.bin, err)
	}
	err = json.Unmarshal(stdoutBuf.Bytes(), stdout)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON output: %w. Please check the command output and try again.", err)
	}
	return nil
}
