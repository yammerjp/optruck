package datasource

import "github.com/joho/godotenv"

type EnvFile struct {
	Path string
}

func NewEnvFile(path string) *EnvFile {
	return &EnvFile{Path: path}
}

func (e *EnvFile) Get() (map[string]string, error) {
	env, err := godotenv.Read(e.Path)
	if err != nil {
		return nil, err
	}
	return env, nil
}
