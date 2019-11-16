package interfaces

type Broker interface {
	Instance() interface{}

	Connect() error
	Publish(string, interface{}) error
	Subscribe(string, interface{}) error
	Close() error
}
