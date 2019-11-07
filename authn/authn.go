package authn

import (
	"errors"
	"time"

	"github.com/advancedlogic/box/commons"
)

type AuthN interface {
	Login(string, string) (interface{}, error)
	Logout(string) error

	Register(string, string) (interface{}, error)
	Delete(string) error
	Reset(string, string) (interface{}, error)
}

type Option func(AuthN) error

type User struct {
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	Timestamp int64    `json:"timestamp"`
	Groups    []string `json:"groups"`
	Enabled   bool     `json:"enabled"`
}

func NewUser(username, password string) (*User, error) {
	if username != "" && password != "" {
		epassword, err := commons.HashAndSalt(password)
		if err != nil {
			return nil, err
		}
		return &User{
			Username:  username,
			Password:  epassword,
			Timestamp: time.Now().UnixNano(),
			Groups:    []string{"user"},
			Enabled:   true,
		}, nil
	}
	return nil, errors.New("username and password cannot be empty")
}
