package interfaces

type Transport interface {
	Instance() interface{}

	Listen() error
	Stop() error

	Get(string, interface{})
	Post(string, interface{})
	Put(string, interface{})
	Delete(string, interface{})
	Static(string, string)
}
