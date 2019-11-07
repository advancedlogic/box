package nats

import (
	"errors"

	"github.com/advancedlogic/box/broker"
	"github.com/advancedlogic/box/logger"
	"github.com/nats-io/go-nats"
)

const (
	errorEndpointEmpty = "endpoint cannot be empty"
	errorCannotCloseConnection = "broker cannot be closed"
)

type Nats struct {
	logger.Logger

	conn          *nats.Conn
	endpoint      string
	handlers      map[string]func(*nats.Msg)
	subscriptions map[string]*nats.Subscription
}

func WithEndpoint(endpoint string) broker.Option {
	return func(i broker.Broker) error {
		if endpoint != "" {
			n := i.(*Nats)
			n.endpoint = endpoint
			return nil
		}
		return errors.New(errorEndpointEmpty)
	}
}

func (n *Nats) Connect() error {
	var conn *nats.Conn
	if conn, err := nats.Connect(n.endpoint); err != nil {
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

func (n *Nats) Publish(topic string, message interface{}) error {
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

func (n *Nats) Close() error {
	if n.conn != nil {
		n.conn.Close()
		return nil
	}
	return errors.New(errorCannotCloseConnection)
}
