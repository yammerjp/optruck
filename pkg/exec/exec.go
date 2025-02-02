package exec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"

	"k8s.io/utils/exec"
)

type Command struct {
	exec   exec.Interface
	logger *slog.Logger
	bin    string
	args   []string
}

func NewCommand(exec exec.Interface, logger *slog.Logger, bin string, args ...string) *Command {
	return &Command{
		exec:   exec,
		logger: logger,
		bin:    bin,
		args:   args,
	}
}

func (c *Command) Run(stdin []byte) (stdout string, err error) {
	execCmd := c.exec.Command(c.bin, c.args...)
	stdoutBytes := bytes.NewBuffer(nil)
	stderrBytes := bytes.NewBuffer(nil)
	if stdin != nil {
		execCmd.SetStdin(bytes.NewReader(stdin))
	}
	execCmd.SetStdout(stdoutBytes)
	execCmd.SetStderr(stderrBytes)
	c.logger.Info(fmt.Sprintf("exec command: %s, %v", c.bin, c.args))
	c.logger.Debug(fmt.Sprintf("stdin: %v", string(stdin)))
	err = execCmd.Run()
	stderr := stderrBytes.String()
	if stderr != "" {
		c.logger.Info(fmt.Sprintf("stderr: %s", stderr))
	}
	if err != nil {
		return "", err
	}
	return stdoutBytes.String(), nil
}

func (c *Command) RunWithJson(stdin interface{}, result interface{}) (err error) {
	stdinBytes, err := json.Marshal(stdin)
	if err != nil {
		return err
	}
	stdout, err := c.Run(stdinBytes)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(stdout), &result)
}
