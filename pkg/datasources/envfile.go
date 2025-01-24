package datasources

import (
	"github.com/joho/godotenv"
)

type SourceType int

const (
	EnvFile SourceType = iota
	KubernetesSecret
)

type Source interface {
	FetchSecrets() (map[string]string, error)
}

type EnvFileSource struct {
	Path string
}

func (e *EnvFileSource) FetchSecrets() (map[string]string, error) {
	return godotenv.Read(e.Path)
}

func NewSource(path string, format SourceType) Source {
	return &EnvFileSource{Path: path}
}
