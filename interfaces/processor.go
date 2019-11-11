package interfaces

//Processor is the "plugin" to specialize the microservice
type Processor interface {
	Init(micro Micro) error
	Close() error
	Process(interface{}) (interface{}, error)
}
