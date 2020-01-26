package nats

import (
	"sync"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

func TestNewNats(t *testing.T) {
	n, _ := New()
	assert.NotEqual(t, n, nil)
}

func TestNats_Connect(t *testing.T) {
	n, _ := New()
	err := n.Connect()
	defer n.Close()
	assert.Equal(t, err, nil)
}

func TestNats_Close(t *testing.T) {
	n, _ := New()
	err := n.Connect()
	assert.Equal(t, err, nil)
	err = n.Close()
	assert.Equal(t, err, nil)
}

func TestNats_PublishSubscribe(t *testing.T) {
	wg := sync.WaitGroup{}
	n, _ := New()
	err := n.Connect()
	assert.Equal(t, err, nil)
	defer n.Close()
	err = n.Subscribe("test", func(msg *nats.Msg) {
		assert.Equal(t, string(msg.Data), []byte("test"))
		wg.Done()
	})
	assert.Equal(t, err, nil)
	err = n.Publish("test", "test")
	assert.Equal(t, err, nil)
	wg.Wait()
}
