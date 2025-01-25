package actions

import (
	"fmt"
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

func (c *MirrorConfig) BuildOpClient() *op.Client {
	return op.NewClient(c.Target)
}

// Mirror は、シークレットをアップロードしてテンプレートを生成するアクションを実行します
func Mirror(config MirrorConfig) error {
	config.Logger.Info("Starting mirror action for item", "item", config.Target.ItemName)

	// データソースからシークレットを取得
	secrets, err := config.DataSource.FetchSecrets()
	if err != nil {
		return fmt.Errorf("failed to fetch secrets from data source: %w", err)
	}
	config.Logger.Info("Fetched secrets from data source", "count", len(secrets))

	secretsResp, err := (*config.BuildOpClient()).UploadItem(secrets, config.Overwrite)
	if err != nil {
		return fmt.Errorf("failed to upload secrets to 1Password: %w", err)
	}
	config.Logger.Info("Uploaded secrets to 1Password successfully")

	err = config.Dest.Write(secretsResp, config.Overwrite)
	if err != nil {
		return fmt.Errorf("failed to write output template: %w", err)
	}
	config.Logger.Info("Template written to %s successfully", "path", config.Dest.GetPath())

	config.Logger.Info("Mirror action completed successfully")
	return nil
}
