package actions

import (
	"log/slog"

	"github.com/yammerjp/optruck/pkg/datasources"
	"github.com/yammerjp/optruck/pkg/op"
	"github.com/yammerjp/optruck/pkg/output"
)

type MirrorConfig struct {
	OpItemClient op.ItemClient
	DataSource   datasources.Source
	Dest         output.Dest
	Overwrite    bool
	Confirmation func() error
}

func (config MirrorConfig) Run() error {
	slog.Debug("Starting mirror action for item", "item", config.OpItemClient.ItemName)

	if err := config.Confirmation(); err != nil {
		slog.Error("failed to confirm", "error", err)
		return err
	}

	secrets, err := config.DataSource.FetchSecrets()
	if err != nil {
		slog.Error("failed to fetch secrets from data source", "error", err)
		return err
	}
	slog.Debug("Fetched secrets from data source", "count", len(secrets))

	secretsResp, err := config.OpItemClient.UploadItem(secrets, config.Overwrite)
	if err != nil {
		slog.Error("failed to upload secrets to 1Password", "error", err)
		return err
	}
	slog.Debug("Uploaded secrets to 1Password successfully")

	err = config.Dest.Write(secretsResp)
	if err != nil {
		slog.Error("failed to write output template", "error", err)
		return err
	}
	slog.Debug("Template written to %s successfully", "path", config.Dest.GetPath())

	slog.Debug("Mirror action completed successfully")
	return nil
}
