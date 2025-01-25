package datasources

import (
	"github.com/joho/godotenv"
)

type EnvFileSource struct {
	Path string
}

func (e *EnvFileSource) FetchSecrets() (map[string]string, error) {
	return godotenv.Read(e.Path)
}
