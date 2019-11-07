package client

type Client interface {
	GET(interface{}) error
	POST(interface{}) error
	PUT(interface{}) error
	DELETE(interface{}) error
}

type Option func(Client) error
