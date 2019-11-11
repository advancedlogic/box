package box

import (
	"errors"
	"io/ioutil"

	"github.com/advancedlogic/box/configuration/viper"
	"github.com/advancedlogic/box/interfaces"
	"github.com/advancedlogic/box/micro"
	"github.com/advancedlogic/box/micro/box"
	"github.com/google/uuid"

	go_shutdown_hook "github.com/ankit-arora/go-utils/go-shutdown-hook"
)

//Box is the main struct for creating a microservice
type Box struct {
	id        string
	name      string
	isRunning bool
	logo      string

	logger        interfaces.Logger
	configuration interfaces.Configuration
	broker        interfaces.Broker
	transport     interfaces.Transport
	client        interfaces.Client
	cache         interfaces.Cache
	registry      interfaces.Registry
	authN         interfaces.AuthN
	authZ         interfaces.AuthZ
	store         interfaces.Store
	processors    []interfaces.Processor
}

type Option func(Box) error

//WithID(id string) set the id of the µs
func WithID(id string) micro.Option {
	return func(m interfaces.Micro) error {
		if id == "" {
			return errors.New("ID cannot be empty")
		}
		box := m.(box.Box)
		box.id = id
		return nil
	}
}

//WithName(name string) set the id of the µs
func WithName(name string) micro.Option {
	return func(m interfaces.Micro) error {
		if name == "" {
			return errors.New("name cannot be empty")
		}
		box := m.(box.Box)
		box.name = name
		return nil
	}
}

func WithLogo(logo interface{}) micro.Option {
	return func(m interfaces.Micro) error {
		if logo != nil {
			box := m.(box.Box)
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

func WithRegistry(registry interfaces.Registry) Option {
	return func(box Box) error {
		if registry != nil {
			box.registry = registry
			return nil
		}
		return errors.New("registry cannot be nil")
	}
}

func WithTransport(transport interfaces.Transport) Option {
	return func(box Box) error {
		if transport != nil {
			box.transport = transport
			return nil
		}
		return errors.New("transport cannot be nil")
	}
}

func WithBroker(broker interfaces.Broker) Option {
	return func(box Box) error {
		if broker != nil {
			box.broker = broker
			return nil
		}
		return errors.New("broker cannot be nil")
	}
}

func WithClient(client interfaces.Client) Option {
	return func(box Box) error {
		if client != nil {
			box.client = client
			return nil
		}
		return errors.New("client cannot be nil")
	}
}

func WithProcessors(processors ...interfaces.Processor) Option {
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

func WithProcessor(processor interfaces.Processor) Option {
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

func WithStore(store interfaces.Store) Option {
	return func(box Box) error {
		if store != nil {
			box.store = store
			return nil
		}
		return errors.New("store cannot be nil")
	}
}

func WithCache(cache interfaces.Cache) Option {
	return func(box Box) error {
		if cache != nil {
			err := cache.Connect()
			if err != nil {
				return err
			}
			box.cache = cache
			return nil
		}
		return errors.New("cache cannot be nil")
	}
}

func WithConfiguration(configuration interfaces.Configuration) Option {
	return func(box Box) error {
		if configuration != nil {
			box.configuration = configuration
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
				viper.WithLogger(box.logger))
			if err != nil {
				return err
			}
			if err := conf.Open(); err != nil {
				return err
			}
			box.configuration = conf
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
				viper.WithLogger(box.logger))
			if err != nil {
				return nil
			}
			box.configuration = conf
		}

		return errors.New("provider and uri cannot be empty")
	}
}

func WithLogger(logger interfaces.Logger) Option {
	return func(box Box) error {
		if logger != nil {
			box.logger = logger
			return nil
		}
		return errors.New("logger cannot be null")
	}
}

func WithAuthN(authn interfaces.AuthN) Option {
	return func(box Box) error {
		if authn != nil {
			box.authN = authn
			return nil
		}
		return errors.New("authn cannot be nil")
	}
}

func WithAuthZ(authz interfaces.AuthZ) Option {
	return func(box Box) error {
		if authz != nil {
			box.authZ = authz
			return nil
		}
		return errors.New("authz cannot be nil")
	}
}

func New(options ...Option) (*Box, error) {
	box := Box{
		id:   uuid.New().String(),
		name: "default",
	}

	for _, option := range options {
		err := option(box)
		if err != nil {
			return nil, err
		}
	}

	return &box, nil
}

func (b *Box) Run() {
	if b.logo != "" {
		println(b.logo)
	}

	go_shutdown_hook.ADD(func() {
		b.Stop()
		b.logger.Warn("Goodbye and thanks for all the fish")
	})
	if b.registry != nil {
		b.logger.Info("registry setup")
		err := b.registry.Register()
		if err != nil {
			b.logger.Fatal(err.Error())
		}
	}
	if b.broker != nil {
		b.logger.Info("broker setup")
		err := b.broker.Connect()
		if err != nil {
			b.logger.Fatal(err.Error())
		}
	}

	if b.transport != nil {
		b.logger.Info("transport setup")
		err := b.transport.Listen()
		if err != nil {
			b.logger.Fatal(err.Error())
		}
	}

	b.isRunning = true
	go_shutdown_hook.Wait()
}

func (b *Box) Stop() {
	if b.broker != nil {
		b.broker.Close()
	}

	if b.transport != nil {
		b.transport.Stop()
	}

	if b.cache != nil {
		b.cache.Close()
	}
}

func (b *Box) Logger() interfaces.Logger {
	return b.logger
}

func (b *Box) Configuration() interfaces.Configuration {
	return b.configuration
}

func (b *Box) Cache() interfaces.Cache {
	return b.cache
}

func (b *Box) Broker() interfaces.Broker {
	return b.broker
}

func (b *Box) Client() interfaces.Client {
	return b.client
}

func (b *Box) Transport() interfaces.Transport {
	return b.transport
}

func (b *Box) Registry() interfaces.Registry {
	return b.registry
}

func (b *Box) AuthN() interfaces.AuthN {
	return b.authN
}

func (b *Box) AuthZ() interfaces.AuthZ {
	return b.AuthZ()
}

func (b *Box) Store() interfaces.Store {
	return b.store
}
