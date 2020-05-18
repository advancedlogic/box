package consul

import (
	"errors"
	"fmt"
	"os"

	"github.com/advancedlogic/box/interfaces"
	"github.com/advancedlogic/box/registry"
	"github.com/hashicorp/consul/api"
)

type Client struct {
	*api.Client
	interfaces.Logger
	address        string
	port           int
	healthEndpoint string
	interval       string
	timeout        string
	username       string
	password       string
}

func WithLogger(logger interfaces.Logger) registry.Option {
	return func(t interfaces.Registry) error {
		if logger != nil {
			client := t.(*Client)
			client.Logger = logger
			return nil
		}
		return errors.New("logger cannot be nil")
	}
}

func WithAddress(address string) registry.Option {
	return func(r interfaces.Registry) error {
		if address != "" {
			client := r.(*Client)
			client.address = address
			return nil
		}

		return errors.New("address cannot be empty")
	}
}

func WithPort(port int) registry.Option {
	return func(r interfaces.Registry) error {
		if port > 0 {
			client := r.(*Client)
			client.port = port
			return nil
		}
		return errors.New("port must be greater than zero")
	}
}

func WithHealthEndpoint(endpoint string) registry.Option {
	return func(i interfaces.Registry) error {
		if endpoint != "" {
			client := i.(*Client)
			client.healthEndpoint = endpoint
			return nil
		}
		return errors.New("endpoint cannot be empty")
	}
}

func WithInterval(interval string) registry.Option {
	return func(i interfaces.Registry) error {
		if interval != "" {
			c := i.(*Client)
			c.interval = interval
			return nil
		}
		return errors.New("interval cannot be empty")
	}
}

func WithTimeout(timeout string) registry.Option {
	return func(i interfaces.Registry) error {
		if timeout != "" {
			c := i.(*Client)
			c.timeout = timeout
			return nil
		}
		return errors.New("timeout cannot be empty")
	}
}

func WithCredentials(username, password string) registry.Option {
	return func(i interfaces.Registry) error {
		if username != "" && password != "" {
			c := i.(*Client)
			c.username = username
			c.password = password
			return nil
		}
		return errors.New("username/password cannot be empty")
	}
}

func New(options ...registry.Option) (*Client, error) {
	client := &Client{
		address:        "localhost:8500",
		username:       "",
		password:       "",
		interval:       "3s",
		timeout:        "5s",
		healthEndpoint: "http://localhost:8080/healthcheck",
	}
	for _, option := range options {
		option(client)
	}
	config := api.DefaultConfig()
	config.Address = client.address
	if client.username != "" && client.password != "" {
		config.HttpAuth = &api.HttpBasicAuth{
			Username: client.username,
			Password: client.password,
		}
	}
	consul, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	client.Client = consul
	return client, nil
}

func (c *Client) Register(name string) error {
	hostname := func() string {
		hn, err := os.Hostname()
		if err != nil {
			c.Fatal(err.Error())
		}
		return hn
	}

	address := hostname()
	registration := &api.AgentServiceRegistration{
		ID:      name,
		Name:    name,
		Port:    c.port,
		Address: address,
	}

	if c.healthEndpoint != "" {
		registration.Check = new(api.AgentServiceCheck)
		registration.Check.HTTP = fmt.Sprintf(c.healthEndpoint)
		registration.Check.Interval = c.interval
		registration.Check.Timeout = c.timeout
	}

	return c.Client.Agent().ServiceRegister(registration)
}

// DeRegister a service with consul local agent
func (c *Client) DeRegister(id string) error {
	return c.Client.Agent().ServiceDeregister(id)
}

// Service return a service
func (c *Client) Service(service, tag string) (interface{}, interface{}, error) {
	passingOnly := true
	addrs, meta, err := c.Client.Health().Service(service, tag, passingOnly, nil)
	if len(addrs) == 0 && err == nil {
		return nil, nil, fmt.Errorf("service ( %s ) was not found", service)
	}
	if err != nil {
		return nil, nil, err
	}
	return addrs, meta, nil
}
