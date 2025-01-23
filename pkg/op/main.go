package op

import (
	"k8s.io/utils/exec"
)

type Client struct {
	exec exec.Interface
}

func NewClient(exec exec.Interface) *Client {
	return &Client{
		exec: exec,
	}
}

func BuildClient() *Client {
	return NewClient(exec.New())
}
