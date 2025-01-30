package actions

import (
	"log/slog"

	"github.com/yammerjp/optruck/pkg/datasources"
	"github.com/yammerjp/optruck/pkg/op"
	"github.com/yammerjp/optruck/pkg/output"
)

type MirrorConfig struct {
	Logger     *slog.Logger
	Target     op.Target
	DataSource datasources.Source
	Dest       output.Dest
	Overwrite  bool
}

func (config MirrorConfig) Run() error {
	config.Logger.Debug("Starting mirror action for item", "item", config.Target.ItemName)

	secrets, err := config.DataSource.FetchSecrets()
	if err != nil {
		config.Logger.Error("failed to fetch secrets from data source", "error", err)
		return err
	}
	config.Logger.Debug("Fetched secrets from data source", "count", len(secrets))

	secretsResp, err := config.Target.BuildClient().UploadItem(secrets, config.Overwrite)
	if err != nil {
		config.Logger.Error("failed to upload secrets to 1Password", "error", err)
		return err
	}
	config.Logger.Debug("Uploaded secrets to 1Password successfully")

	err = config.Dest.Write(secretsResp)
	if err != nil {
		config.Logger.Error("failed to write output template", "error", err)
		return err
	}
	config.Logger.Debug("Template written to %s successfully", "path", config.Dest.GetPath())

	config.Logger.Debug("Mirror action completed successfully")
	return nil
}
