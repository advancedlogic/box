package interfaces

type Registry interface {
	Register(string) error
	DeRegister(string) error
	Service(string, string) (interface{}, interface{}, error)
}
