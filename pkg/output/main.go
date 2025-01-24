package output

type Client struct {
	Format              Format
	EnvTemplateFilePath string
}

type Format string

const (
	FormatEnvTemplate   Format = "env"
	FormatK8sSecretYaml Format = "k8s-secret"
)

func NewClient(format string, envTemplateFilePath string) *Client {
	return &Client{Format: Format(format), EnvTemplateFilePath: envTemplateFilePath}
}
