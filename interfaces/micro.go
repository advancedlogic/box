package interfaces

type Micro interface {
	Run()
	Stop()

	Configuration() Configuration
	Registry() Registry
	Broker() Broker
	Transport() Transport
	Logger() Logger
	Processor() Processor
	AuthN() AuthN
	AuthZ() AuthZ
	Cache() Cache
	Client() Client
	Store() Store
}
