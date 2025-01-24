package output

import (
	"errors"

	"github.com/yammerjp/optruck/pkg/op"
)

type K8sSecretDest struct {
	Path string
}

func (d *K8sSecretDest) GetPath() string {
	return d.Path
}

func (d *K8sSecretDest) Write(resp *op.SecretResponse) error {
	return errors.New("not implemented")
}
