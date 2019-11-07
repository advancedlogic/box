package fs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/advancedlogic/box/authn"
	"github.com/advancedlogic/box/commons"
)

type FS struct {
	folder string
}

func WithFolder(folder string) authn.Option {
	return func(a authn.AuthN) error {
		if folder != "" {
			fs := a.(*FS)
			fs.folder = folder
			return nil
		}
		return errors.New("folder cannot be empty")
	}
}

func New(options ...authn.Option) (*FS, error) {
	fs := &FS{
		folder: "fs",
	}
	for _, option := range options {
		if err := option(fs); err != nil {
			return nil, err
		}
	}
	return fs, nil
}

func (f *FS) Register(username, password string) (interface{}, error) {
	if username != "" && password != "" {
		user, err := authn.NewUser(username, password)
		if err != nil {
			return nil, err
		}
		jsonUser, err := json.Marshal(user)
		if err != nil {
			return nil, err
		}
		err = ioutil.WriteFile(fmt.Sprintf("%s/%s.json", f.folder, username), jsonUser, 0644)
		if err != nil {
			return nil, err
		}
		user.Password = ""
		return user, nil
	}
	return nil, errors.New("username and password cannot be empty")
}

func (f *FS) Login(username, password string) (interface{}, error) {
	if username != "" && password != "" {
		jsonUser, err := ioutil.ReadFile(fmt.Sprintf("%s/%s.json", f.folder, username))
		if err != nil {
			return nil, err
		}
		var user authn.User
		err = json.Unmarshal(jsonUser, &user)
		if err != nil {
			return nil, err
		}

		if !commons.ComparePasswords(user.Password, []byte(password)) {
			return nil, errors.New("wrong username or password")
		}
		user.Password = ""
		return user, nil
	}
	return nil, errors.New("username and password cannot be empty")
}

func (f *FS) Logout(username string) error {
	if username != "" {
		return nil
	}
	return errors.New("username cannot be empty")
}

func (f *FS) Delete(username string) error {
	if username != "" {
		return nil
	}
	return errors.New("username cannot be empty")
}

func (f *FS) Reset(username, password string) (interface{}, error) {
	return f.Register(username, password)
}
