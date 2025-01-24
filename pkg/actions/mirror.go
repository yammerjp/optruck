package actions

import (
	"fmt"
	"log"

	"github.com/yammerjp/optruck/pkg/datasources"
	"github.com/yammerjp/optruck/pkg/op"
	"github.com/yammerjp/optruck/pkg/output"
)

// MirrorConfig は、mirror アクションに必要な設定を表します
type MirrorConfig struct {
	ItemName     string             // 1Passwordのアイテム名
	Vault        string             // 1PasswordのVault名
	DataSource   datasources.Source // データソース（.envファイルやKubernetes Secretなど）
	OutputPath   string             // 出力ファイルのパス
	OutputFormat string             // 出力フォーマット ("env" または "k8s")
	Overwrite    bool               // 上書きオプション
}

// Mirror は、シークレットをアップロードしてテンプレートを生成するアクションを実行します
func Mirror(config MirrorConfig) error {
	log.Printf("Starting mirror action for item: %s", config.ItemName)

	// データソースからシークレットを取得
	secrets, err := config.DataSource.FetchSecrets()
	if err != nil {
		return fmt.Errorf("failed to fetch secrets from data source: %w", err)
	}
	log.Printf("Fetched %d secrets from data source", len(secrets))

	// 1Passwordにシークレットをアップロード
	opClient := op.NewClient("", config.Vault) // 必要に応じてアカウント情報を渡す
	resp, err := opClient.CreateItem(config.ItemName, secrets)
	if err != nil {
		return fmt.Errorf("failed to upload secrets to 1Password: %w", err)
	}
	log.Println("Uploaded secrets to 1Password successfully")

	outputClient := output.NewClient(config.OutputFormat, config.OutputPath)
	err = outputClient.Write(resp, config.ItemName)
	if err != nil {
		return fmt.Errorf("failed to write output template: %w", err)
	}
	log.Printf("Template written to %s successfully", config.OutputPath)

	log.Println("Mirror action completed successfully")
	return nil
}
