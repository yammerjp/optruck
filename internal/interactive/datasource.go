package interactive

import (
	"github.com/manifoldco/promptui"
)

type DataSourceEnum int

const (
	DataSourceEnvFile DataSourceEnum = iota
	DataSourceK8sSecret
)

func (ds DataSourceEnum) String() string {
	switch ds {
	case DataSourceEnvFile:
		return "env file"
	case DataSourceK8sSecret:
		return "k8s secret"
	default:
		return "unknown"
	}
}

func (r Runner) SelectDataSource() (DataSourceEnum, error) {
	dataSources := []struct {
		Label string
		Value DataSourceEnum
	}{
		{Label: "env file", Value: DataSourceEnvFile},
		{Label: "k8s secret", Value: DataSourceK8sSecret},
	}
	i, _, err := r.Select(promptui.Select{
		Label:     "Select data source: ",
		Items:     dataSources,
		Templates: SelectTemplateBuilder("Data Source", "Label", ""),
	})
	if err != nil {
		return 0, err
	}
	return dataSources[i].Value, nil
}
