package output

import (
	"fmt"
	"os"
	"text/template"

	"github.com/yammerjp/optruck/pkg/op"
)

func (d *EnvTemplateDest) GetPath() string {
	return d.Path
}

func (d *EnvTemplateDest) Write(resp *op.SecretResponse) error {
	file, err := os.OpenFile(d.Path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open env template file: %v", err)
	}
	defer file.Close()

	tmpl, err := template.New("envtemplate").Parse(`
#   - 1password vault: {{.VaultName}}
{{if .AccountName}}
#   - 1password account: {{.AccountName}}
{{end}}
# To restore, run the following command:
#   $ cat .env.1password | grep -v '^#' | op inject > .env

{{range .FieldLabels}}
{{.Label}}={{"{{op://"}}{{ .VaultName }}/{{ .ItemName }}/{{ .Label }}{{"}}"}}
{{end}}
`)
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	return tmpl.Execute(file, resp)
}
