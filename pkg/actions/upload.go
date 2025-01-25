package actions

import (
	"fmt"
	"log/slog"

	"github.com/yammerjp/optruck/pkg/datasources"
	"github.com/yammerjp/optruck/pkg/op"
	"github.com/yammerjp/optruck/pkg/output"
)

type UploadConfig struct {
	Logger     *slog.Logger
	Target     op.Target
	DataSource datasources.Source // データソース（.envファイルやKubernetes Secretなど）
	Dest       output.Dest        // 出力先
	Overwrite  bool               // 上書きオプション
}

func (config UploadConfig) Run() error {
	config.Logger.Info("Starting upload action for item", "item", config.Target.ItemName)

	secrets, err := config.DataSource.FetchSecrets()
	if err != nil {
		return fmt.Errorf("failed to fetch secrets from data source: %w", err)
	}
	config.Logger.Info("Fetched secrets from data source", "count", len(secrets))

	_, err = config.Target.BuildClient().UploadItem(secrets, config.Overwrite)
	if err != nil {
		return fmt.Errorf("failed to upload secrets to 1Password: %w", err)
	}
	config.Logger.Info("Uploaded secrets to 1Password successfully")
	return nil
}
