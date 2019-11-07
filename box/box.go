package box

import (
	"github.com/advancedlogic/box/transport"
	"github.com/advancedlogic/box/broker"
	"github.com/advancedlogic/box/configuration"
	"github.com/advancedlogic/box/logger"
)

//Box is the main struct for creating a microservice
type Box struct {
	logger.Logger
	configuration.Configuration
	broker.Broker
	transport.Transport
}
