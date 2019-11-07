package configuration

//Configuration is the base configuration interface
type Configuration interface {
	Open(...string) error
	Get(string) (interface{}, error)
}

//Option is an helper to inject parameters during the initialization
type Option func(Configuration) error