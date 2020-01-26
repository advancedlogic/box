package nats

import (
	"errors"
	"reflect"

	"github.com/advancedlogic/box/broker"
	"github.com/advancedlogic/box/interfaces"
	"github.com/nats-io/nats.go"
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

func New(options ...broker.Option) (*Nats, error) {
	nats := &Nats{
		endpoint:      "localhost:4222",
		handlers:      make(map[string]func(*nats.Msg)),
		subscriptions: make(map[string]*nats.Subscription),
	}
	for _, option := range options {
		if err := option(nats); err != nil {
			return nil, err
		}
	}
	return nats, nil
}

func (n *Nats) Instance() interface{} {
	return n.conn
}

func (n *Nats) Connect() error {
	conn, err := nats.Connect(n.endpoint)
	if err != nil {
		return err
	}
	n.conn = conn
	for topic, handler := range n.handlers {
		subscription, err := n.conn.QueueSubscribe(topic, "default", handler)
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
	println(reflect.TypeOf(handler))
	f := handler.(func(msg *nats.Msg))
	n.handlers[topic] = f
	return nil
}

func (n Nats) Close() error {
	if n.conn != nil {
		n.conn.Close()
		return nil
	}
	return errors.New(errorCannotCloseConnection)
}
