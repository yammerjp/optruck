package task

import (
	"fmt"
	"log/slog"

	"github.com/yammerjp/optruck/pkg/datasource"
	"github.com/yammerjp/optruck/pkg/op"
	"github.com/yammerjp/optruck/pkg/output"
)

type MirrorTask struct {
	logger *slog.Logger

	itemName            string
	envFilePath         string
	envTemplateFilePath string
}

func NewMirrorTask(logger *slog.Logger, itemName string, envFilePath string, envTemplateFilePath string) *MirrorTask {
	return &MirrorTask{
		logger:              logger,
		itemName:            itemName,
		envFilePath:         envFilePath,
		envTemplateFilePath: envTemplateFilePath,
	}
}

func (t *MirrorTask) Run() error {
	t.logger.Info("mirroring secrets")
	t.logger.Info(fmt.Sprintf("itemName: %s", t.itemName))
	t.logger.Info(fmt.Sprintf("envFilePath: %s", t.envFilePath))

	env, err := datasource.NewEnvFile(t.envFilePath).Get()
	if err != nil {
		return err
	}
	t.logger.Debug(fmt.Sprintf("env: %v", env))

	opClient := op.BuildClient()
	resp, err := opClient.CreateItem("", "", t.itemName, env)
	if err != nil {
		return err
	}
	t.logger.Debug(fmt.Sprintf("resp: %v", resp))

	outputClient := output.NewClient(t.logger, t.envTemplateFilePath)
	outputClient.Print(resp)
	t.logger.Debug(fmt.Sprintf("envTemplateFilePath: %s", t.envTemplateFilePath))

	return nil

}

func RunMirror(logger *slog.Logger, itemName string, envFilePath string, envTemplateFilePath string) error {
	task := NewMirrorTask(logger, itemName, envFilePath, envTemplateFilePath)
	return task.Run()
}
