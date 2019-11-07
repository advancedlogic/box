package resty

import (
	"crypto/tls"
	"errors"
	"net/http"

	"github.com/advancedlogic/box/client"
	"gopkg.in/resty.v1"
)

type Resty struct {
	Url         string
	QueryParams map[string]string
	Headers     map[string]string
	Cookies     map[string]string
	AuthToken   string
	Body        string
	Username    string
	Password    string
	pem         string
	key         string
}

func WithUrl(url string) client.Option {
	return func(client client.Client) error {
		if url != "" {
			r := client.(*Resty)
			r.Url = url
			return nil
		}
		return errors.New("url cannot be empty")
	}
}

func AddQueryParam(key, value string) client.Option {
	return func(client client.Client) error {
		if key != "" && value != "" {
			r := client.(*Resty)
			r.QueryParams[key] = value
			return nil
		}
		return errors.New("key and value cannot be empty")
	}
}

func AddHeader(key, value string) client.Option {
	return func(client client.Client) error {
		if key != "" && value != "" {
			r := client.(*Resty)
			r.Headers[key] = value
			return nil
		}
		return errors.New("key and value cannot be empty")
	}
}

func AddCookie(key, value string) client.Option {
	return func(client client.Client) error {
		if key != "" && value != "" {
			r := client.(*Resty)
			r.Cookies[key] = value
			return nil
		}
		return errors.New("key and value cannot be empty")
	}
}

func WithAuthToken(token string) client.Option {
	return func(client client.Client) error {
		if token != "" {
			r := client.(*Resty)
			r.AuthToken = token
			return nil
		}
		return errors.New("token cannot be empty")
	}
}

func WithBody(body string) client.Option {
	return func(client client.Client) error {
		if body != "" {
			r := client.(*Resty)
			r.Body = body
			return nil
		}
		return errors.New("body cannot be empty")
	}
}

func WithBasicAuthentication(username, password string) client.Option {
	return func(client client.Client) error {
		if username != "" && password != "" {
			r := client.(*Resty)
			r.Username = username
			r.Password = password
			return nil
		}
		return errors.New("username and password cannot be empty")
	}
}

func WithX509Certificate(pem, key string) client.Option {
	return func(client client.Client) error {
		if pem != "" && key != "" {
			r := client.(*Resty)
			r.key = key
			r.pem = pem
			return nil
		}
		return errors.New("pem and key cannot be empty")
	}
}

func New(options ...client.Option) (*Resty, error) {
	r := &Resty{
		QueryParams: make(map[string]string),
		Headers:     make(map[string]string),
		Cookies:     make(map[string]string),
	}
	for _, option := range options {
		if err := option(r); err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (r *Resty) render() (*resty.Request, error) {
	client := resty.New()
	if len(r.Cookies) > 0 {
		for key, value := range r.Cookies {
			client.SetCookie(&http.Cookie{
				Name:  key,
				Value: value,
			})
		}
	}

	if r.pem != "" && r.key != "" {
		cert, err := tls.LoadX509KeyPair(r.pem, r.key)
		if err != nil {
			return nil, err
		}
		client.SetCertificates(cert)
	}

	request := client.R()
	if len(r.Headers) > 0 {
		request.SetHeaders(r.Headers)
	}
	if len(r.QueryParams) > 0 {
		request.SetQueryParams(r.QueryParams)
	}

	if r.AuthToken != "" {
		request.SetAuthToken(r.AuthToken)
	}
	if r.Body != "" {
		request.SetBody(r.Body)
	}
	if r.Username != "" && r.Password != "" {
		request.SetBasicAuth(r.Username, r.Password)
	}
	return request, nil
}

func (r *Resty) GET(h interface{}) error {
	request, err := r.render()
	if err != nil {
		return err
	}
	response, err := request.Get(r.Url)
	if err != nil {
		return err
	}
	handler := h.(func(response *resty.Response) error)
	return handler(response)
}

func (r *Resty) POST(h interface{}) error {
	request, err := r.render()
	if err != nil {
		return err
	}
	response, err := request.Post(r.Url)
	if err != nil {
		return err
	}
	handler := h.(func(response *resty.Response) error)
	return handler(response)
}

func (r *Resty) PUT(h interface{}) error {
	request, err := r.render()
	if err != nil {
		return err
	}
	response, err := request.Put(r.Url)
	if err != nil {
		return err
	}
	handler := h.(func(response *resty.Response) error)
	return handler(response)
}

func (r *Resty) DELETE(h interface{}) error {
	request, err := r.render()
	if err != nil {
		return err
	}
	response, err := request.Delete(r.Url)
	if err != nil {
		return err
	}
	handler := h.(func(response *resty.Response) error)
	return handler(response)
}
