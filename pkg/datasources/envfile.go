package datasources

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/yammerjp/optruck/internal/errors"
)

type EnvFileSource struct {
	Path string
}

func (e *EnvFileSource) FetchSecrets() (map[string]string, error) {
	envPairs, err := godotenv.Read(e.Path)
	if err != nil {
		return nil, errors.NewOperationFailedError(
			fmt.Sprintf("環境変数ファイル '%s' の読み込み", e.Path),
			err,
			"ファイルが存在し、読み取り権限があることを確認してください",
		)
	}

	if len(envPairs) == 0 {
		return nil, errors.NewInvalidArgumentError(
			"環境変数ファイル",
			"環境変数が設定されていません",
			"環境変数ファイルに値が設定されていることを確認してください",
		)
	}

	return envPairs, nil
}
