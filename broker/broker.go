package broker

type Broker interface {
	Connect() error
	Publish(string, interface{}) error
	Subscribe(string, interface{}) error
	Close() error
}

type Option func(Broker) error
