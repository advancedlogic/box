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

func (v Viper) Instance() interface{} {
	return v.Viper
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
func (v *Viper) Get(key string) interface{} {
	return v.Viper.Get(key)
}

func (v *Viper) Default(key string, def interface{}) interface{} {
	value := v.Get(key)
	if value == nil {
		return def
	}
	return value
}

func (v *Viper) String(key string, def string) string {
	value := v.GetString(key)
	if value == "" {
		return def
	}
	return value
}

func (v *Viper) Int(key string, def int) int {
	value := v.GetInt(key)
	if value == 0 {
		return def
	}
	return value
}

func (v *Viper) Int32(key string, def int32) int32 {
	value := v.GetInt32(key)
	if value == 0 {
		return def
	}
	return value
}

func (v *Viper) Int64(key string, def int64) int64 {
	value := v.GetInt64(key)
	if value == 0 {
		return def
	}
	return value
}

func (v *Viper) Float(key string, def float64) float64 {
	value := v.GetFloat64(key)
	if value == 0.0 {
		return def
	}
	return value
}

func (v *Viper) Bool(key string, def bool) bool {
	value := v.GetBool(key)
	if value == false {
		return def
	}
	return value
}

func (v *Viper) MapOfStrings(path string, def map[string]string) map[string]string {
	if value := v.Viper.GetStringMapString(path); value != nil {
		return value
	}
	return def
}

func (v *Viper) ArrayOfStrings(path string, def []string) []string {
	if value := v.GetStringSlice(path); value != nil {
		return value
	}
	return def
}
