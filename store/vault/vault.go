package vault

import (
	"fmt"
	"time"

	"github.com/advancedlogic/box/commons"
	"github.com/advancedlogic/box/store"
	"github.com/advancedlogic/box/interfaces"
	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
)

type Vault struct {
	id                  string
	namespace           string
	token               string
	servers             []string
	timeout             time.Duration
	skipTLSVerification bool
}

func WithToken(token string) store.Option {
	return func(s interfaces.Store) error {
		if token != "" {
			vault := s.(*Vault)
			vault.token = token
			return nil
		}
		return errors.New("token cannot be empty")
	}
}

func WithNamespace(namespace string) store.Option {
	return func(s interfaces.Store) error {
		if namespace != "" {
			vault := s.(*Vault)
			vault.namespace = namespace
			return nil
		}
		return errors.New("token cannot be empty")
	}
}

func WithServers(servers ...string) store.Option {
	return func(s interfaces.Store) error {
		if len(servers) > 0 {
			vault := s.(*Vault)
			for _, server := range servers {
				vault.servers = append(vault.servers, server)
			}
			return nil
		}
		return errors.New("at least one server must be provided")
	}
}

func SkipTLSVerification(skip bool) store.Option {
	return func(s interfaces.Store) error {
		vault := s.(*Vault)
		vault.skipTLSVerification = skip
		return nil
	}
}

func New(options ...store.Option) (*Vault, error) {
	v := &Vault{
		id:                  commons.UUID(),
		namespace:           "default",
		token:               "",
		servers:             make([]string, 0),
		timeout:             10 * time.Second,
		skipTLSVerification: true,
	}

	for _, option := range options {
		if err := option(v); err != nil {
			return nil, err
		}
	}

	return v, nil
}

func (v *Vault) connect() (*api.Client, error) {
	config := &api.Config{
		Address: v.servers[0],
	}
	if err := config.ConfigureTLS(&api.TLSConfig{
		Insecure: v.skipTLSVerification,
	}); err != nil {
		return nil, err
	}
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	client.SetToken(v.token)
	return client, nil
}

func (v *Vault) Create(key string, value interface{}) error {
	client, err := v.connect()
	if err != nil {
		return err
	}

	_, err = client.Logical().Write(fmt.Sprintf("/%s/%s", v.namespace, key), value.(map[string]interface{}))
	if err != nil {
		return err
	}
	return nil
}

func (v *Vault) Read(key string) (interface{}, error) {
	client, err := v.connect()
	if err != nil {
		return nil, err
	}

	secret, err := client.Logical().Read(fmt.Sprintf("/%s/%s", v.namespace, key))
	if err != nil {
		return nil, err
	}

	return secret.Data, nil
}

func (v *Vault) Update(key string, value interface{}) error {
	return v.Create(key, value)
}

func (v *Vault) Delete(key string) error {
	client, err := v.connect()
	if err != nil {
		return err
	}

	_, err = client.Logical().Delete(fmt.Sprintf("/%s/%s", v.namespace, key))
	if err != nil {
		return err
	}
	return nil
}

func (v *Vault) List(params ...interface{}) (interface{}, error) {
	client, err := v.connect()
	if err != nil {
		return nil, err
	}
	secret, err := client.  //.Logical().List(fmt.Sprintf("/%s", v.namespace))
	if err != nil {
		return nil, err
	}
	return secret.Data, nil
}

func (v *Vault) Query(params ...interface{}) (interface{}, error) {
	return nil, nil
}
