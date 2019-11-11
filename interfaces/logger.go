package interfaces

//Logger defines the interface for logging the application
type Logger interface {
	Instance() interface{}

	Info(string)
	Debug(string)
	Warn(string)
	Error(string)
	Fatal(string)
}

