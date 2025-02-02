package interactiverunner

import (
	"github.com/manifoldco/promptui"
)

type DataSourceEnum int

const (
	DataSourceEnvFile DataSourceEnum = iota
	DataSourceK8sSecret
)

type DataSourceSelector struct {
	runner InteractiveRunner
}

func NewDataSourceSelector(runner InteractiveRunner) *DataSourceSelector {
	return &DataSourceSelector{runner: runner}
}

func (ds *DataSourceSelector) Select() (DataSourceEnum, error) {
	dataSources := []struct {
		Label string
		Value DataSourceEnum
	}{
		{Label: "env file", Value: DataSourceEnvFile},
		{Label: "k8s secret", Value: DataSourceK8sSecret},
	}
	i, _, err := ds.runner.Select(promptui.Select{
		Label:     "Select data source: ",
		Items:     dataSources,
		Templates: SelectTemplateBuilder("Data Source", "Label", ""),
	})
	if err != nil {
		return 0, err
	}
	return dataSources[i].Value, nil
}
