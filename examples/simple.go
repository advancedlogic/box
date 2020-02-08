package main

import (
	"net/http"

	"github.com/advancedlogic/box/box"
	"github.com/advancedlogic/box/configuration/viper"
	"github.com/advancedlogic/box/interfaces"
	"github.com/advancedlogic/box/logger/logrus"
	"github.com/advancedlogic/box/transport/rest"
	"github.com/advancedlogic/box/registry/consul"
	"github.com/gin-gonic/gin"
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

	transport, err := rest.New(
		rest.WithPort(9999),
		rest.WithLogger(logger),
	)
	if err != nil {
		panic(err)
	}

	transport.Get("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "It works")
	})

	registry, err := consul.New(
		consul.WithLogger(logger),
		consul.WithPort(9999),
	)
	if err != nil {
		panic(err)
	}

	processor := Simple{}

	box, err := box.New(
		box.WithName("simple"),
		box.WithLogo("banner.txt"),
		box.WithLogger(logger),
		box.WithConfiguration(configuration),
		box.WithProcessor(processor),
		box.WithTransport(transport),
		box.WithRegistry(registry),
	)
	if err != nil {
		panic(err)
	}

	box.Run()
}
