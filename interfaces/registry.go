package interfaces

type Registry interface {
	Instance() interface{}
	Register() error
	Search(string) (interface{}, error)
}
