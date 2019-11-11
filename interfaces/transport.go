package interfaces

type Transport interface {
	Listen() error
	Stop() error

	Get(string, interface{})
	Post(string, interface{})
	Put(string, interface{})
	Delete(string, interface{})
}
