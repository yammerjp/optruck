package output

import "github.com/yammerjp/optruck/pkg/op"

type Dest interface {
	Write(resp *op.SecretReference, overwrite bool) error
	GetPath() string
	GetBasename() string
}

var _ Dest = (*EnvTemplateDest)(nil)
var _ Dest = (*K8sSecretTemplateDest)(nil)
