package actions

import (
	"log/slog"

	"github.com/yammerjp/optruck/pkg/datasources"
	"github.com/yammerjp/optruck/pkg/op"
	"github.com/yammerjp/optruck/pkg/output"
)

// MirrorConfig は、mirror アクションに必要な設定を表します
type MirrorConfig struct {
	Logger     *slog.Logger
	Target     op.Target
	DataSource datasources.Source // データソース（.envファイルやKubernetes Secretなど）
	Dest       output.Dest        // 出力先
	Overwrite  bool               // 上書きオプション
}

// Mirror は、シークレットをアップロードしてテンプレートを生成するアクションを実行します
func (config MirrorConfig) Run() error {
	config.Logger.Debug("Starting mirror action for item", "item", config.Target.ItemName)

	// データソースからシークレットを取得
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

	err = config.Dest.Write(secretsResp, config.Overwrite)
	if err != nil {
		config.Logger.Error("failed to write output template", "error", err)
		return err
	}
	config.Logger.Debug("Template written to %s successfully", "path", config.Dest.GetPath())

	config.Logger.Debug("Mirror action completed successfully")
	return nil
}
