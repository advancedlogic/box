package store

type Store interface {
	Create(string, interface{}) error
	Read(string) (interface{}, error)
	Update(string, interface{}) error
	Delete(string) error
	List(...interface{}) (interface{}, error)
	Query(...interface{}) (interface{}, error)
}

type Option func(Store) error
