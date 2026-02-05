package keyring

import (
	"github.com/99designs/keyring"
)

type Keyring struct {
	ring keyring.Keyring
}

func NewRing() (*Keyring, error) {
	ring, err := keyring.Open(keyring.Config{
		ServiceName: "papyro",
	})
	if err != nil {
		return &Keyring{}, err
	}

	return &Keyring{ring: ring}, nil
}

func (kr *Keyring) Get(name string) (keyring.Item, error) {
	return kr.ring.Get(name)
}

func (kr *Keyring) Save(name string, password string) error {
	return kr.ring.Set(keyring.Item{
		Key:  name,
		Data: []byte(password),
	})
}

func (kr *Keyring) Delete(name string) error {
	return kr.ring.Remove(name)
}
