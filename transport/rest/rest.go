package rest

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/advancedlogic/box/interfaces"
	"github.com/advancedlogic/box/transport"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/toorop/gin-logrus"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

type Rest struct {
	interfaces.Logger
	port         int
	health       string
	readTimeout  time.Duration
	writeTimeout time.Duration
	cert         string
	key          string
	server       *http.Server
	www          string
	router       *gin.Engine
}

func WithLogger(logger interfaces.Logger) transport.Option {
	return func(t interfaces.Transport) error {
		if logger != nil {
			rest := t.(*Rest)
			rest.Logger = logger
			return nil
		}
		return errors.New("logger cannot be nil")
	}
}

func WithPort(port int) transport.Option {
	return func(t interfaces.Transport) error {
		if port > 0 {
			rest := t.(*Rest)
			rest.port = port
			return nil
		}
		return errors.New("port cannot be zero")
	}
}

func WithHealthCheckEndpoint(health string) transport.Option {
	return func(t interfaces.Transport) error {
		if health != "" {
			rest := t.(*Rest)
			rest.health = health
			return nil
		}
		return errors.New("health check endpoint cannot be empty")
	}
}

func WithReadTimeout(timeout time.Duration) transport.Option {
	return func(i interfaces.Transport) error {
		if timeout != 0 {
			r := i.(*Rest)
			r.readTimeout = timeout
		}
		return errors.New("timeout cannot be zero")
	}
}

func WithWriteTimeout(timeout time.Duration) transport.Option {
	return func(i interfaces.Transport) error {
		if timeout != 0 {
			r := i.(*Rest)
			r.writeTimeout = timeout
		}
		return errors.New("timeout cannot be zero")
	}
}

func (r *Rest) scanPort(ip string, port int, timeout time.Duration) error {
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(timeout)
			err = r.scanPort(ip, port, timeout)
		}
		return err
	}

	if err = conn.Close(); err != nil {
		return err
	}
	r.Warn(fmt.Sprintf("port %s is busy", port))
	return nil
}

func (r *Rest) findAlternativePort() error {
	currentPort := r.port
	for port := currentPort; port < 32000; port++ {
		err := r.scanPort("localhost", port, 10*time.Second)
		if err != nil {
			r.port = port
			return nil
		}
	}
	return errors.New("no alternatives port found")
}

func New(options ...transport.Option) (*Rest, error) {
	rest := &Rest{
		port:         8080,
		readTimeout:  5 * time.Second,
		writeTimeout: 5 * time.Second,
		router:       gin.New(),
	}

	for _, option := range options {
		if err := option(rest); err != nil {
			return nil, err
		}
	}

	return rest, nil
}

func (r *Rest) Instance() interface{} {
	return r.router
}

func (r *Rest) Listen() error {
	router := r.router
	logger := r.Instance().(*logrus.Logger)
	router.Use(ginlogrus.Logger(logger), gin.Recovery())
	router.GET(r.health, func(c *gin.Context) {
		c.String(200, "transport service is good")
	})

	p := ginprometheus.NewPrometheus("gin")
	p.Use(router)

	if err := r.findAlternativePort(); err != nil {
		r.Fatal(err.Error())
	}

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", r.port),
		Handler:        router,
		ReadTimeout:    r.readTimeout,
		WriteTimeout:   r.writeTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	r.server = s
	httpHeader := "http"
	go func() {
		if r.cert != "" && r.key != "" {
			if err := s.ListenAndServeTLS(r.cert, r.key); err != nil {
				r.Fatal(err.Error())
			}
			httpHeader += "s"
		} else if err := s.ListenAndServe(); err != nil {
			r.Fatal(err.Error())
		}

	}()
	r.Info(fmt.Sprintf("Http(s) server listening on port %d", r.port))
	return nil
}

func (r *Rest) Stop() error {
	return nil
}

func (r *Rest) Get(url string, h interface{}) {
	handler := h.(func(c *gin.Context))
	r.router.GET(url, handler)
}

func (r *Rest) Post(url string, h interface{}) {
	handler := h.(func(c *gin.Context))
	r.router.POST(url, handler)
}

func (r *Rest) Put(url string, h interface{}) {
	handler := h.(func(c *gin.Context))
	r.router.PUT(url, handler)
}

func (r *Rest) Delete(url string, h interface{}) {
	handler := h.(func(c *gin.Context))
	r.router.DELETE(url, handler)
}

func (r *Rest) Static(url string, folder string) {
	r.router.Static(url, folder)
}
