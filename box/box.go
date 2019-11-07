package box

import (
	"errors"
	"io/ioutil"

	"github.com/advancedlogic/box/authn"
	"github.com/advancedlogic/box/authz"
	"github.com/advancedlogic/box/broker"
	"github.com/advancedlogic/box/cache"
	"github.com/advancedlogic/box/client"
	"github.com/advancedlogic/box/configuration"
	"github.com/advancedlogic/box/configuration/viper"
	"github.com/advancedlogic/box/logger"
	"github.com/advancedlogic/box/processor"
	"github.com/advancedlogic/box/registry"
	"github.com/advancedlogic/box/store"
	"github.com/advancedlogic/box/transport"
)

//Box is the main struct for creating a microservice
type Box struct {
	id        string
	name      string
	isRunning bool
	logo      string

	logger.Logger
	configuration.Configuration
	broker.Broker
	transport.Transport
	client.Client
	cache.Cache
	registry.Registry
	authn.AuthN
	authz.AuthZ
	store.Store
	processors []processor.Processor
}

type Option func(Box) error

//WithID(id string) set the id of the µs
func WithID(id string) Option {
	return func(box Box) error {
		if id == "" {
			return errors.New("ID cannot be empty")
		}
		box.id = id
		return nil
	}
}

//WithName(name string) set the id of the µs
func WithName(name string) Option {
	return func(box Box) error {
		if name == "" {
			return errors.New("name cannot be empty")
		}
		box.name = name
		return nil
	}
}

func WithLogo(logo interface{}) Option {
	return func(box Box) error {
		if logo != nil {
			switch logo.(type) {
			case []byte:
				box.logo = string(logo.([]byte))
			case string:
				b, err := ioutil.ReadFile(logo.(string))
				if err != nil {
					return err
				}
				box.logo = string(b)
			}
			return nil
		}
		return errors.New("logo cannot be nil")
	}
}

func WithRegistry(registry registry.Registry) Option {
	return func(box Box) error {
		if registry != nil {
			box.Registry = registry
			return nil
		}
		return errors.New("registry cannot be nil")
	}
}

func WithTransport(transport transport.Transport) Option {
	return func(box Box) error {
		if transport != nil {
			box.Transport = transport
			return nil
		}
		return errors.New("transport cannot be nil")
	}
}

func WithBroker(broker broker.Broker) Option {
	return func(box Box) error {
		if broker != nil {
			box.Broker = broker
			return nil
		}
		return errors.New("broker cannot be nil")
	}
}

func WithClient(client client.Client) Option {
	return func(box Box) error {
		if client != nil {
			box.Client = client
			return nil
		}
		return errors.New("client cannot be nil")
	}
}

func WithProcessors(processors ...processor.Processor) Option {
	return func(box Box) error {
		if processors != nil && len(processors) > 0 {
			for _, processor := range processors {
				err := processor.Init(box)
				if err != nil {
					return err
				}
				box.processors = append(box.processors, processor)
				return nil
			}
		}
		return errors.New("processors cannot be nil or empty")
	}
}

func WithProcessor(processor processor.Processor) Option {
	return func(box Box) error {
		if processor != nil {
			err := processor.Init(box)
			if err != nil {
				return err
			}
			box.processors = append(box.processors, processor)
			return nil
		}
		return errors.New("processor cannot be nil")
	}
}

func WithStore(store store.Store) Option {
	return func(box Box) error {
		if store != nil {
			box.Store = store
			return nil
		}
		return errors.New("store cannot be nil")
	}
}

func WithCache(cache cache.Cache) Option {
	return func(box Box) error {
		if cache != nil {
			err := cache.Init()
			if err != nil {
				return err
			}
			box.Cache = cache
			return nil
		}
		return errors.New("cache cannot be nil")
	}
}

func WithConfiguration(configuration configuration.Configuration) Option {
	return func(box Box) error {
		if configuration != nil {
			box.Configuration = configuration
			return nil
		}
		return errors.New("configuration cannot be nil")
	}
}

func WithLocalConfiguration() Option {
	return func(box Box) error {
		if box.name != "" {
			conf, err := viper.New(
				viper.WithName(box.name),
				viper.WithLogger(box.Logger))
			if err != nil {
				return err
			}
			if err := conf.Open(); err != nil {
				return err
			}
			box.Configuration = conf
		}
		return errors.New("name cannot be empty")
	}
}

func WithRemoteConfiguration(provider, uri string) Option {
	return func(box Box) error {
		if provider != "" && uri != "" {
			conf, err := viper.New(
				viper.WithName(box.name),
				viper.WithProvider(provider),
				viper.WithURI(uri),
				viper.WithLogger(box.Logger))
			if err != nil {
				return nil
			}
			box.Configuration = conf
		}

		return errors.New("provider and uri cannot be empty")
	}
}

func WithLogger(logger logger.Logger) Option {
	return func(box Box) error {
		if logger != nil {
			box.Logger = logger
			return nil
		}
		return errors.New("logger cannot be null")
	}
}

func WithAuthN(authn authn.AuthN) Option {
	return func(box Box) error {
		if authn != nil {
			box.AuthN = authn
			return nil
		}
		return errors.New("authn cannot be nil")
	}
}

func WithAuthZ(authz authn.AuthZ) Option {
	return func(box Box) error {
		if authz != nil {
			box.AuthN = authz
			return nil
		}
		return errors.New("authz cannot be nil")
	}
}
