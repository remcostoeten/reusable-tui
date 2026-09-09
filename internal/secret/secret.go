package secret

import (
	"errors"

	"github.com/zalando/go-keyring"
)

var ErrNotFound = errors.New("secret: not found")

type Vault interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Delete(key string) error
}

type Keyring struct {
	Service string
}

func NewKeyring(service string) Keyring {
	return Keyring{Service: service}
}

func (k Keyring) Get(key string) (string, error) {
	value, err := keyring.Get(k.Service, key)
	if errors.Is(err, keyring.ErrNotFound) {
		return "", ErrNotFound
	}
	return value, err
}

func (k Keyring) Set(key, value string) error {
	return keyring.Set(k.Service, key, value)
}

func (k Keyring) Delete(key string) error {
	err := keyring.Delete(k.Service, key)
	if errors.Is(err, keyring.ErrNotFound) {
		return ErrNotFound
	}
	return err
}
