package cache

type Cache interface {
	Instance() interface{}

	Connect() error
	Close() error
	
	Set(string, interface{}, int) error
	Get(string) (interface{}, error)
	Keys() (interface{}, error)
}