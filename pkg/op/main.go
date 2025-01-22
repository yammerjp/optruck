package op

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func BuildClient() *Client {
	return NewClient()
}
