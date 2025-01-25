package datasources

type Source interface {
	FetchSecrets() (map[string]string, error)
}
