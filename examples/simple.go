package main

import (
	"github.com/advancedlogic/box/box"
	"github.com/advancedlogic/box/configuration/viper"
	"github.com/advancedlogic/box/interfaces"
	"github.com/advancedlogic/box/logger/logrus"
)

type Simple struct {
	interfaces.Micro
}

func (s Simple) Init(micro interfaces.Micro) error {
	s.Micro = micro
	return nil
}

func (s Simple) Process(data interface{}) (interface{}, error) {
	s.Logger().Info("This is a simple example")
	return nil, nil
}

func (s Simple) Close() error {
	return nil
}

func main() {
	configuration, err := viper.New()
	if err != nil {
		panic(err)
	}

	logger, err := logrus.New(logrus.WithLevel(configuration.Default("level", "info").(string)))
	if err != nil {
		panic(err)
	}

	processor := Simple{}

	box, err := box.New(
		box.WithLogger(logger),
		box.WithConfiguration(configuration),
		box.WithProcessor(processor),
	)
	if err != nil {
		panic(err)
	}

	box.Run()
}
