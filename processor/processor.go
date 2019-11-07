package processor

import "github.com/advancedlogic/box/box"

//Processor is the "plugin" to specialize the microservice
type Processor interface {
	Init(box box.Box) error
	Close() error
	Process(interface{}) (interface{}, error)
}

type Option func(Processor) error