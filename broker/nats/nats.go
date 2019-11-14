package nats

import (
	"errors"

	"github.com/advancedlogic/box/broker"
	"github.com/advancedlogic/box/interfaces"
	nats "github.com/nats-io/nats.go"
)

const (
	errorEndpointEmpty         = "endpoint cannot be empty"
	errorCannotCloseConnection = "broker cannot be closed"
	errorLoggerNil             = "logger cannot be nil"
)

type Nats struct {
	interfaces.Logger

	conn          *nats.Conn
	endpoint      string
	handlers      map[string]func(*nats.Msg)
	subscriptions map[string]*nats.Subscription
}

func WithEndpoint(endpoint string) broker.Option {
	return func(i interfaces.Broker) error {
		if endpoint != "" {
			n := i.(*Nats)
			n.endpoint = endpoint
			return nil
		}
		return errors.New(errorEndpointEmpty)
	}
}

func WithLogger(logger interfaces.Logger) broker.Option {
	return func(i interfaces.Broker) error {
		if logger != nil {
			n := i.(*Nats)
			n.Logger = logger
			return nil
		}
		return errors.New(errorLoggerNil)
	}
}

func (n *Nats) Connect() error {
	conn, err := nats.Connect(n.endpoint)
	if err != nil {
		return err
	}
	n.conn = conn
	for topic, handler := range n.handlers {
		subscription, err := n.conn.Subscribe(topic, handler)
		if err != nil {
			return err
		}
		n.subscriptions[topic] = subscription
	}
	return nil
}

func (n Nats) Publish(topic string, message interface{}) error {
	var m []byte
	switch message.(type) {
	case string:
		m = []byte(message.(string))
	default:
		m = message.([]byte)
	}
	return n.conn.Publish(topic, m)
}

func (n *Nats) Subscribe(topic string, handler interface{}) error {
	n.handlers[topic] = handler.(func(*nats.Msg))
	return nil
}

func (n Nats) Close() error {
	if n.conn != nil {
		n.conn.Close()
		return nil
	}
	return errors.New(errorCannotCloseConnection)
}
