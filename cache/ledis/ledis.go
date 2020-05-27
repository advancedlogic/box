package ledis

import (
	"context"
	"errors"

	"github.com/advancedlogic/box/cache"
	"github.com/advancedlogic/box/interfaces"
	"github.com/go-redis/redis/v8"
)

type Ledis struct {
	collection    string
	endpoints     []string
	username      string
	password      string
	db            int
	clusterClient *redis.ClusterClient
	client        *redis.Client
	ctx           context.Context
}

func WithCollection(collection string) cache.Option {
	return func(c interfaces.Cache) error {
		if collection != "" {
			ledis := c.(*Ledis)
			ledis.collection = collection
			return nil
		}
		return errors.New("collection cannot be empty")
	}
}

func WithPassowrd(password string) cache.Option {
	return func(c interfaces.Cache) error {
		if password != "" {
			ledis := c.(*Ledis)
			ledis.password = password
			return nil
		}
		return errors.New("username and password cannot be empty")
	}
}

func WithDB(db int) cache.Option {
	return func(c interfaces.Cache) error {
		if db > -1 {
			ledis := c.(*Ledis)
			ledis.db = db
			return nil
		}
		return errors.New("db must be >= 0")
	}
}

func AddEndpoints(endpoints ...string) cache.Option {
	return func(c interfaces.Cache) error {
		for _, endpoint := range endpoints {
			if endpoint != "" {
				ledis := c.(*Ledis)
				ledis.endpoints = append(ledis.endpoints, endpoint)
				return nil
			}
		}
		return errors.New("endpoint cannot be empty")
	}
}

func New(options ...cache.Option) (*Ledis, error) {
	ledis := &Ledis{
		endpoints: make([]string, 0),
		ctx:       context.Background(),
	}
	for _, option := range options {
		if err := option(ledis); err != nil {
			return nil, err
		}
	}
	return ledis, nil
}

func (l *Ledis) Instance() interface{} {
	if l.client != nil {
		return l.client
	}
	return l.clusterClient
}

func (l *Ledis) Connect() error {
	if len(l.endpoints) > 1 {
		clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    l.endpoints,
			Password: l.password,
		})
		_, err := clusterClient.Ping(l.ctx).Result()
		if err != nil {
			return err
		}
		l.clusterClient = clusterClient
	} else {
		client := redis.NewClient(&redis.Options{
			Addr:     l.endpoints[0],
			Password: l.password,
			DB:       l.db,
		})
		l.client = client
	}
	return nil
}

func (l *Ledis) Close() error {
	if l.clusterClient != nil {
		return l.clusterClient.Close()
	}
	return nil
}

func (l *Ledis) Set(key string, value interface{}, ttl int) error {
	var status *redis.StatusCmd
	if l.client != nil {
		status = l.client.Set(l.ctx, key, value, -1)
	} else {
		status = l.clusterClient.Set(l.ctx, key, value, -1)
	}
	if status.Err() != nil {
		return status.Err()
	}
	return nil
}

func (l *Ledis) Get(key string) (interface{}, error) {
	var status *redis.StringCmd
	if l.client != nil {
		status = l.client.Get(l.ctx, key)
	} else {
		status = l.clusterClient.Get(l.ctx, key)
	}
	if status.Err() != nil {
		return nil, status.Err()
	}
	result := status.Val()
	return result, nil
}

func (l *Ledis) Keys() (interface{}, error) {
	var status *redis.StringSliceCmd
	if l.client != nil {
		status = l.client.Keys(l.ctx, "*")
	} else {
		status = l.clusterClient.Keys(l.ctx, "*")
	}
	if status.Err() != nil {
		return nil, status.Err()
	}
	return status.Val(), nil
}
