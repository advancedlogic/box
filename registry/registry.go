package registry

type Registry interface {
	Register() error
}

type Option func(Registry) error