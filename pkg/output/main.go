package output

import "github.com/yammerjp/optruck/pkg/op"

type Format string

const (
	EnvFile   Format = "env"
	K8sSecret Format = "k8s"
)

type Dest interface {
	Write(resp *op.SecretResponse) error
	GetPath() string
}

type EnvTemplateDest struct {
	Path string
}

func NewDest(path string, format Format) Dest {
	if format == EnvFile {
		return &EnvTemplateDest{Path: path}
	} else if format == K8sSecret {
		return &K8sSecretDest{Path: path}
	}
	return nil
}
