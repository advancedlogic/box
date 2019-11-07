package cache

type Cache interface {
	Connect() error
	Close() error
	Set(string, interface{}, int) error
	Get(string) (interface{}, error)
}