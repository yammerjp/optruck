package output

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/yammerjp/optruck/pkg/op"
)

type K8sSecretDest struct {
	Path       string
	Namespace  string
	SecretName string
}

func (d *K8sSecretDest) GetPath() string {
	return d.Path
}

type k8sTemplateData struct {
	*op.SecretResponse
	*K8sSecretDest
}

func (d *K8sSecretDest) GetBasename() string {
	return filepath.Base(d.Path)
}

func (d *K8sSecretDest) Write(resp *op.SecretResponse) error {
	file, err := os.Create(d.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	tmpl, err := template.New("k8s-secret").Parse(`# This file was generated by optruck.{{if .SecretResponse.Account}}
#   - 1password account: {{.SecretResponse.Account}}{{end}}{{if .SecretResponse.VaultName}}
#   - 1password vault: {{.SecretResponse.VaultName}}{{end}}
# To restore, run the following command:
#   $ op inject -i {{.K8sSecretDest.GetBasename}} | kubectl apply -f -

apiVersion: v1
kind: Secret
metadata:
  name: {{.K8sSecretDest.SecretName}}
  namespace: {{.K8sSecretDest.Namespace}}
type: Opaque
data:{{range .SecretResponse.GetFieldRefs}}
  {{.Label}}: {{.Ref}}{{end}}`)
	if err != nil {
		return err
	}

	return tmpl.Execute(file, k8sTemplateData{
		SecretResponse: resp,
		K8sSecretDest:  d,
	})
}
