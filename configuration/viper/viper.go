package viper

import (
	"errors"
	"fmt"

	"github.com/advancedlogic/box/interfaces"
	"github.com/fsnotify/fsnotify"

	"github.com/advancedlogic/box/configuration"
	"github.com/spf13/viper"
)

//WithName specifies the name of the configuration to open
func WithName(name string) configuration.Option {
	return func(i interfaces.Configuration) error {
		if name != "" {
			v := i.(*Viper)
			v.name = name
			return nil
		}

		return errors.New("provider cannot be empty")
	}
}

//WithProvider specifies the name of the provider in case
//you want to manage a remote configuration (e.g. consul, etcd)
func WithProvider(provider string) configuration.Option {
	return func(i interfaces.Configuration) error {
		if provider != "" {
			v := i.(*Viper)
			v.provider = provider
			return nil
		}

		return errors.New("provider cannot be empty")
	}
}

//WithURI specifies the address of the provider
func WithURI(uri string) configuration.Option {
	return func(i interfaces.Configuration) error {
		if uri != "" {
			v := i.(*Viper)
			v.uri = uri
			return nil
		}

		return errors.New("uri cannot be empty")
	}
}

//Viper is a wrapper around the viper library
type Viper struct {
	*viper.Viper

	name     string
	provider string
	uri      string
}

//New create a new configuration based on the given options
func New(options ...configuration.Option) (*Viper, error) {
	v := &Viper{
		Viper: viper.New(),
	}
	for _, option := range options {
		if err := option(v); err != nil {
			return nil, err
		}
	}

	return v, nil
}

//Open one or more configuration files
func (v Viper) Open(paths ...string) error {
	v.SetConfigName(v.name)

	if v.provider != "" && v.uri != "" {
		if err := v.AddRemoteProvider(v.provider, v.uri, v.name); err != nil {
			return err
		}
		if err := v.ReadRemoteConfig(); err != nil {
			return err
		}
	} else {
		v.AddConfigPath(fmt.Sprintf("/etc/%s/", v.name))
		v.AddConfigPath(fmt.Sprintf("$HOME/.%s", v.name))
		for _, path := range paths {
			v.AddConfigPath(path)
		}
		v.AddConfigPath(".")
		v.AutomaticEnv()
		if err := v.ReadInConfig(); err != nil {
			return err
		}
	}

	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		err := v.ReadInConfig()
		if err != nil {
			return
		}
	})

	return nil
}

//Get return a configuration property given a key
func (v *Viper) Get(key string) (interface{}, error) {
	return v.Get(key)
}

func (v *Viper) Default(key string, def interface{}) interface{} {
	value, err := v.Get(key)
	if err != nil {
		return def
	} 

	return value
}
