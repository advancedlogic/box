package logrus

import (
	"errors"

	"github.com/advancedlogic/box/logger"
	"github.com/sirupsen/logrus"
)

const (
	errorLevelEmpty  = "Logger level cannot be empty. Use info, warn, error, fatal or debug"
	errorFormatEmpty = "Logger format cannot be empty. Check for logrus specs for formatting"
)

//WithLevel received a level as string.
//Possible values are: info, warn, error, fatal, debug.
//Default value is info.
func WithLevel(level string) logger.Option {
	return func(i logger.Logger) error {
		if level != "" {
			l := i.(Logrus)
			l.level = level
			return nil
		}
		return errors.New(errorLevelEmpty)
	}
}

//WithFormat received a format as a string
//For more details about logrus format specs check their website
func WithFormat(format string) logger.Option {
	return func(i logger.Logger) error {
		if format != "" {
			l := i.(Logrus)
			l.format = format
			return nil
		}
		return errors.New(errorFormatEmpty)
	}
}

//Logrus is a struct implementing the Logger interface
//Basically is a wrapper around the logrus library
type Logrus struct {
	*logrus.Logger

	level  string
	format string
}

//New instantiate a new Logger with the given options
func New(options ...logger.Option) (*Logrus, error) {
	l := &Logrus{}
	for _, option := range options {
		err := option(l)
		if err != nil {
			return nil, err
		}
	}

	//Default values if options are empty
	if l.level == "" {
		l.level = "info"
	}

	switch l.level {
	case "info":
		l.SetLevel(logrus.InfoLevel)
	case "warn":
		l.SetLevel(logrus.WarnLevel)
	case "error":
		l.SetLevel(logrus.ErrorLevel)
	case "fatal":
		l.SetLevel(logrus.FatalLevel)
	case "debug":
		l.SetLevel(logrus.DebugLevel)
	default:
		l.SetLevel(logrus.InfoLevel)
	}

	if l.format == "" {
		l.format = ""
	}

	return l, nil
}

//Instance get the instance of the 
func (l Logrus) Instance() interface{} {
	return l.Logger
}

//Info logging level
func (l Logrus) Info(message string) {
	l.Info(message)
}

//Debug logging level
func (l Logrus) Debug(message string) {
	l.Debug(message)
}

//Warn logging level
func (l Logrus) Warn(message string) {
	l.Warn(message)
}

//Error logging level
func (l Logrus) Error(message string) {
	l.Error(message)
}

//Fatal logging level
func (l Logrus) Fatal(message string) {
	l.Fatal(message)
}
