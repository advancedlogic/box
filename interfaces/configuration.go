package interfaces

//Configuration is the base configuration interface
type Configuration interface {
	Open(...string) error
	Get(string) interface{}
	Default(string, interface{}) interface{}
}
