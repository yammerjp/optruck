package datasources

import (
	"log/slog"

	"github.com/joho/godotenv"
)

type EnvFileSource struct {
	Path   string
	Logger *slog.Logger
}

func (e *EnvFileSource) FetchSecrets() (map[string]string, error) {
	e.Logger.Debug("Reading secrets from env file", "path", e.Path)

	secrets, err := godotenv.Read(e.Path)
	if err != nil {
		e.Logger.Error("Failed to read env file", "path", e.Path, "error", err)
		return make(map[string]string), err
	}

	e.Logger.Debug("Successfully read secrets from env file", "path", e.Path, "count", len(secrets))
	return secrets, nil
}
