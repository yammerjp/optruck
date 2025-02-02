package exec

import (
	"bytes"
	"encoding/json"
	"log/slog"

	"k8s.io/utils/exec"
)

type CommandConfig struct {
	exec   exec.Interface
	logger *slog.Logger
}

func NewCommandConfig(exec exec.Interface, logger *slog.Logger) CommandConfig {
	return CommandConfig{
		exec:   exec,
		logger: logger,
	}
}

type Command struct {
	CommandConfig
	bin  string
	args []string
}

func (cc *CommandConfig) Command(bin string, args ...string) *Command {
	return &Command{
		CommandConfig: *cc,
		bin:           bin,
		args:          args,
	}
}

func NewCommand(exec exec.Interface, logger *slog.Logger, bin string, args ...string) *Command {
	return &Command{
		CommandConfig: CommandConfig{
			exec:   exec,
			logger: logger,
		},
		bin:  bin,
		args: args,
	}
}

func (c *Command) Run(stdin []byte) (stdout string, err error) {
	c.logger.Debug("Executing command", "command", c.bin, "args", c.args)

	execCmd := c.exec.Command(c.bin, c.args...)
	stdoutBytes := bytes.NewBuffer(nil)
	stderrBytes := bytes.NewBuffer(nil)
	if stdin != nil {
		execCmd.SetStdin(bytes.NewReader(stdin))
		c.logger.Debug("Command input", "stdin", string(stdin))
	}
	execCmd.SetStdout(stdoutBytes)
	execCmd.SetStderr(stderrBytes)

	err = execCmd.Run()
	stderr := stderrBytes.String()
	stdout = stdoutBytes.String()

	if stderr != "" {
		c.logger.Info("Command produced stderr output", "stderr", stderr)
	}
	if err != nil {
		c.logger.Error("Command execution failed", "error", err, "stderr", stderr)
		return "", err
	}

	c.logger.Debug("Command executed successfully", "stdout_length", len(stdout))
	return stdout, nil
}

func (c *Command) RunWithJson(stdin interface{}, result interface{}) (err error) {
	c.logger.Debug("Executing JSON command", "command", c.bin, "args", c.args)

	stdinBytes, err := json.Marshal(stdin)
	if err != nil {
		c.logger.Error("Failed to marshal JSON input", "error", err)
		return err
	}

	stdout, err := c.Run(stdinBytes)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(stdout), &result)
	if err != nil {
		c.logger.Error("Failed to unmarshal JSON output", "error", err, "stdout", stdout)
		return err
	}

	c.logger.Debug("JSON command completed successfully")
	return nil
}
