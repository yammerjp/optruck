package datasources

type Source interface {
	FetchSecrets() (map[string]string, error)
}

var _ Source = (*EnvFileSource)(nil)
var _ Source = (*K8sSecretSource)(nil)
