package op

import "k8s.io/utils/exec"

type Target struct {
	Account  string
	Vault    string
	ItemName string
}

type Client struct {
	exec exec.Interface
	Target
}

func (target Target) BuildClient() *Client {
	return &Client{
		exec:   exec.New(),
		Target: target,
	}
}
