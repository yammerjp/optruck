package task

import (
	"fmt"
	"log/slog"

	"github.com/yammerjp/optruck/pkg/datasource"
	"github.com/yammerjp/optruck/pkg/op"
	"github.com/yammerjp/optruck/pkg/output"
)

type MirrorTask struct {
	Logger              *slog.Logger
	AccountName         string
	VaultName           string
	ItemName            string
	EnvFilePath         string
	EnvTemplateFilePath string
}

func (t *MirrorTask) Run() error {
	t.Logger.Debug("mirroring secrets")
	t.Logger.Debug(fmt.Sprintf("ItemName: %s", t.ItemName))
	t.Logger.Debug(fmt.Sprintf("EnvFilePath: %s", t.EnvFilePath))

	env, err := datasource.NewEnvFile(t.EnvFilePath).Get()
	if err != nil {
		return err
	}
	t.Logger.Debug(fmt.Sprintf("Env: %v", env))

	opClient := op.BuildClient()
	resp, err := opClient.CreateItem(t.AccountName, t.VaultName, t.ItemName, env)
	if err != nil {
		return err
	}
	t.Logger.Debug(fmt.Sprintf("resp: %v", resp))

	outputClient := &output.Client{
		EnvTemplateFilePath: t.EnvTemplateFilePath,
		AccountID:           t.AccountName,
		VaultID:             resp.Vault.ID,
		ItemID:              resp.ID,
		ItemName:            resp.Title,
	}
	outputClient.Print(resp)
	t.Logger.Debug(fmt.Sprintf("envTemplateFilePath: %s", t.EnvTemplateFilePath))

	return nil

}
