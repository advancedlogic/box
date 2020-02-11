package minio

import (
	"errors"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/advancedlogic/box/interfaces"
	"github.com/advancedlogic/box/store"
	minio "github.com/minio/minio-go/v6"
)

type Minio struct {
	location  string
	endpoint  string
	bucket    string
	accessKey string
	secretKey string
	lock      sync.RWMutex
	client    minio.Client
}

func WithLocation(location string) store.Option {
	return func(i interfaces.Store) error {
		if location != "" {
			m := i.(*Minio)
			m.location = location
			return nil
		}
		return errors.New("location cannot be empty")
	}
}

func WithBucket(bucket string) store.Option {
	return func(i interfaces.Store) error {
		if bucket != "" {
			m := i.(*Minio)
			m.bucket = bucket
			return nil
		}
		return errors.New("location cannot be empty")
	}
}

func WithEndpoint(endpoint string) store.Option {
	return func(i interfaces.Store) error {
		if endpoint != "" {
			m := i.(*Minio)
			m.endpoint = endpoint
			return nil
		}
		return errors.New("endpoint cannot be empty")
	}
}

func WithAccessKey(accessKey string) store.Option {
	return func(i interfaces.Store) error {
		if accessKey != "" {
			m := i.(*Minio)
			m.accessKey = accessKey
			return nil
		}
		return errors.New("access key cannot be empty")
	}
}

func WithSecretKey(secretKey string) store.Option {
	return func(i interfaces.Store) error {
		if secretKey != "" {
			m := i.(*Minio)
			m.secretKey = secretKey
			return nil
		}
		return errors.New("secret key cannot be empty")
	}
}

func WithCredentials(accessKey, secretKey string) store.Option {
	return func(i interfaces.Store) error {
		if accessKey != "" && secretKey != "" {
			m := i.(*Minio)
			m.accessKey = accessKey
			m.secretKey = secretKey
			return nil
		}
		return errors.New("access or secret key cannot be empty")
	}
}

func New(options ...store.Option) (*Minio, error) {
	m := &Minio{
		location: "default",
		bucket:   "default",
		endpoint: "localhost:9000",
		lock:     sync.RWMutex{},
	}

	for _, option := range options {
		if err := option(m); err != nil {
			return nil, err
		}
	}

	client, err := minio.New(m.endpoint, m.accessKey, m.secretKey, false)
	if err != nil {
		return nil, err
	}
	m.client = client

	return m, nil
}

func (m Minio) Instance() interface{} {
	return nil
}

func (m *Minio) Create(bucket string, key string, data interface{}) error {
	client := m.client
	reader := strings.NewReader(data.(string))

	m.lock.Lock()
	defer m.lock.Unlock()
	_, err = client.PutObject(bucket, key, reader, -1, minio.PutObjectOptions{
		ContentType: "plain/txt",
	})
	if err != nil {
		return err
	}

	return nil
}

func (m *Minio) Read(bucket string, key string) (interface{}, error) {
	client := m.client

	m.lock.Lock()
	defer m.lock.Unlock()
	reader, err := client.GetObject(bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return "", err
	}
	defer reader.Close()

	if value, err := ioutil.ReadAll(reader); err == nil {
		return string(value), nil
	} else {
		return nil, err
	}
}

func (m *Minio) Update(bucket string, key string, data interface{}) error {
	return m.Create(bucket, key, data)
}

func (m *Minio) Delete(bucket string, key string) error {
	client = m.client
	m.lock.Lock()
	defer m.lock.Unlock()
	return client.RemoveObject(bucket, key)
}

func (m *Minio) List(bucket string, params ...interface{}) (interface{}, error) {
	client := m.client
	doneCh := make(chan struct{})
	defer close(doneCh)
	prefix := params[0].(string)
	callback := params[1].(func([]byte))
	for value := range client.ListObjectsV2(bucket, prefix, true, doneCh) {
		reader, err := client.GetObject(bucket, value.Key, minio.GetObjectOptions{})
		if err != nil {
			return "", err
		}

		value, err := ioutil.ReadAll(reader)
		if err == nil {
			callback(value)
		}
		reader.Close()
	}
	return nil, nil
}

func (m *Minio) Query(bucket string, params ...interface{}) (interface{}, error) {
	return nil, nil
}
