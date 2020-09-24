package interfaces

type Store interface {
	Create(string, string, interface{}) error
	Read(string, string) (interface{}, error)
	Update(string, string, interface{}) error
	Delete(string, string) error
	List(string, ...interface{}) (interface{}, error)
	Query(string, ...interface{}) (interface{}, error)
	Buckets() (interface{}, error)
}
